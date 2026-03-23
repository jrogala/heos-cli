package client

import "strings"

// ParseMessage parses a HEOS message field string like "pid=1&level=50"
// into a map of key-value pairs.
func ParseMessage(msg string) map[string]string {
	result := make(map[string]string)
	if msg == "" {
		return result
	}
	pairs := strings.Split(msg, "&")
	for _, pair := range pairs {
		k, v, _ := strings.Cut(pair, "=")
		result[k] = v
	}
	return result
}

// EncodeHEOS encodes special characters for HEOS protocol values.
// '&' -> %26, '=' -> %3D, '%' -> %25
func EncodeHEOS(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "&", "%26")
	s = strings.ReplaceAll(s, "=", "%3D")
	return s
}

// DecodeHEOS decodes HEOS protocol encoded values.
// %26 -> '&', %3D -> '=', %25 -> '%'
func DecodeHEOS(s string) string {
	s = strings.ReplaceAll(s, "%26", "&")
	s = strings.ReplaceAll(s, "%3D", "=")
	s = strings.ReplaceAll(s, "%25", "%")
	return s
}
