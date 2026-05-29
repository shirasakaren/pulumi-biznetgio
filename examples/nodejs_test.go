//go:build nodejs || all
// +build nodejs all

package examples

import "testing"

func TestNodejsExampleLifecycle(t *testing.T) {
	t.Skip("requires pulumi CLI and generated SDKs; re-enable once resources land")
}
