//go:build dotnet || all
// +build dotnet all

package examples

import "testing"

func TestDotnet(t *testing.T) {
	t.Skip("requires pulumi CLI and generated SDKs; re-enable once resources land")
}
// wip 442
