//go:build go || all
// +build go all

package examples

import "testing"

func TestGoExampleLifecycle(t *testing.T) {
	t.Skip("requires pulumi CLI and generated SDKs; re-enable once resources land")
}
// wip 615
