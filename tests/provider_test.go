package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shirasakaren/pulumi-biznetgio/provider"
)

func TestProviderName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "biznetgio", provider.Name)
}
