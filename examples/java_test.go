//go:build java || all
// +build java all

package examples

import "testing"

func TestJava(t *testing.T) {
	t.Skip("java sdk not generated yet; re-enable once resources land")
}
