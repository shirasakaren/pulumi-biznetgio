package provider

import (
	"encoding/json"
	"strings"
)

// secretJSONKeys are keys whose values must be masked before entering state/log.
var secretJSONKeys = []string{
	"cipassword", "console_password", "consolepassword", "password", "passwd",
	"private_key", "privatekey", "private", "secret_key", "secretkey", "secret",
	"pem", "token",
}

func isSecretJSONKey(k string) bool {
	for _, s := range secretJSONKeys {
		if strings.EqualFold(k, s) {
			return true
		}
	}
	return false
}

// redactMap salin map, mask key rahasia (case-insensitive), rekursif ke nested map.
func redactMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if isSecretJSONKey(k) {
			out[k] = "***"
			continue
		}
		if nested, ok := v.(map[string]any); ok {
			out[k] = redactMap(nested)
			continue
		}
		out[k] = v
	}
	return out
}

// RedactJSON marshals a map with secret values masked, falling back to nil.
func RedactJSON(m map[string]any) []byte {
	b, err := json.Marshal(redactMap(m))
	if err != nil {
		return nil
	}
	return b
}
