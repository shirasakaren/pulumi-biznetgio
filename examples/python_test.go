//go:build python || all
// +build python all

package examples

import "testing"

func TestPython(t *testing.T) {
	t.Skip("requires pulumi CLI and generated SDKs; re-enable once resources land")
}
// wip 1
