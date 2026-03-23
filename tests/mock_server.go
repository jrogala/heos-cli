package tests

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
)

// MockResponse holds a preconfigured response for a HEOS command.
type MockResponse struct {
	Result  string          `json:"result"`
	Message string          `json:"message"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// MockServer is a fake HEOS TCP server for testing.
type MockServer struct {
	listener  net.Listener
	mu        sync.Mutex
	responses map[string]MockResponse // "group/command" -> response
	received  []string                // commands received
	done      chan struct{}
}

// NewMockServer creates and starts a mock HEOS server on a random port.
func NewMockServer() (*MockServer, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	s := &MockServer{
		listener:  ln,
		responses: make(map[string]MockResponse),
		done:      make(chan struct{}),
	}
	go s.serve()
	return s, nil
}

// Addr returns the server address (host:port).
func (s *MockServer) Addr() string {
	return s.listener.Addr().String()
}

// Host returns just the host part.
func (s *MockServer) Host() string {
	host, _, _ := net.SplitHostPort(s.Addr())
	return host
}

// Port returns just the port number.
func (s *MockServer) Port() int {
	addr := s.listener.Addr().(*net.TCPAddr)
	return addr.Port
}

// Close shuts down the mock server.
func (s *MockServer) Close() {
	close(s.done)
	s.listener.Close()
}

// On registers a response for a given HEOS command path (e.g. "player/get_players").
func (s *MockServer) On(command string, resp MockResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responses[command] = resp
}

// OnSuccess registers a success response with a message string.
func (s *MockServer) OnSuccess(command, message string) {
	s.On(command, MockResponse{Result: "success", Message: message})
}

// OnSuccessWithPayload registers a success response with a JSON payload.
func (s *MockServer) OnSuccessWithPayload(command, message string, payload any) {
	data, _ := json.Marshal(payload)
	s.On(command, MockResponse{Result: "success", Message: message, Payload: data})
}

// OnFail registers a failure response.
func (s *MockServer) OnFail(command string, eid int, text string) {
	s.On(command, MockResponse{
		Result:  "fail",
		Message: fmt.Sprintf("eid=%d&text=%s", eid, text),
	})
}

// Received returns all commands received by the server.
func (s *MockServer) Received() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, len(s.received))
	copy(out, s.received)
	return out
}

// Reset clears all registered responses and received commands.
func (s *MockServer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responses = make(map[string]MockResponse)
	s.received = nil
}

func (s *MockServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				continue
			}
		}
		go s.handleConn(conn)
	}
}

func (s *MockServer) handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		line := strings.TrimRight(string(buf[:n]), "\r\n")
		cmd := s.parseCommand(line)

		s.mu.Lock()
		s.received = append(s.received, line)
		resp, ok := s.responses[cmd]
		s.mu.Unlock()

		var reply map[string]any
		if ok {
			reply = map[string]any{
				"heos": map[string]any{
					"command": cmd,
					"result":  resp.Result,
					"message": resp.Message,
				},
			}
			if resp.Payload != nil {
				reply["payload"] = json.RawMessage(resp.Payload)
			}
		} else {
			reply = map[string]any{
				"heos": map[string]any{
					"command": cmd,
					"result":  "success",
					"message": "",
				},
			}
		}

		data, _ := json.Marshal(reply)
		conn.Write(append(data, '\r', '\n'))
	}
}

// parseCommand extracts "group/command" from "heos://group/command?params".
func (s *MockServer) parseCommand(line string) string {
	line = strings.TrimPrefix(line, "heos://")
	if idx := strings.Index(line, "?"); idx >= 0 {
		line = line[:idx]
	}
	return line
}
