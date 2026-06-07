//go:build yaml || all
// +build yaml all

package examples

import "testing"

func TestYAMLExampleLifecycle(t *testing.T) {
	t.Skip("requires pulumi CLI and generated SDKs; re-enable once resources land")
}
