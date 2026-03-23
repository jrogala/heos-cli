package tests

import (
	"github.com/jrogala/heos-cli/client"
)

type scenarioCtx struct {
	server *MockServer
	client *client.Client

	// Result holders
	lastErr     error
	players     []client.Player
	player      *client.Player
	groups      []client.Group
	group       *client.Group
	nowPlaying  *client.NowPlaying
	queueItems  []client.QueueItem
	sources     []client.MusicSource
	quickSelects []client.QuickSelect
	account     map[string]string
	stringVal   string
	repeatVal   string
	shuffleVal  string
}

func newScenarioCtx() *scenarioCtx {
	sc := &scenarioCtx{
		server: globalServer,
	}
	sc.client = client.New(globalServer.Host(), globalServer.Port())
	return sc
}

func (sc *scenarioCtx) connect() error {
	return sc.client.Connect()
}

func (sc *scenarioCtx) cleanup() {
	if sc.client != nil {
		sc.client.Close()
	}
}
