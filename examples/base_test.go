package examples

import (
	"github.com/pulumi/providertest/providers"
	goprovider "github.com/pulumi/pulumi-go-provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/biznetgio/pulumi-biznetgio/provider"
)

var providerFactory = func(_ providers.PulumiTest) (pulumirpc.ResourceProviderServer, error) { //nolint:unused
	return goprovider.RawServer("biznetgio", "1.0.0", provider.Provider())(nil)
}
