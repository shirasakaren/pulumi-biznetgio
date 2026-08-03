package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type GpuKeypair struct{}

type GpuKeypairArgs struct {
	Name string `pulumi:"name"`
}

type GpuKeypairState struct {
	GpuKeypairArgs
	PublicKey  *string `pulumi:"publicKey"`
	PrivateKey *string `pulumi:"privateKey" provider:"secret"`
	Raw        string  `pulumi:"raw" provider:"secret"`
}

func (a *GpuKeypairArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Name of the keypair. Create-only, changing it replaces the keypair.")
}

func (s *GpuKeypairState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.PublicKey, "Generated public key.")
	ann.Describe(&s.PrivateKey, "Generated private key. Only returned at creation.")
	ann.Describe(&s.Raw, "Raw JSON of the keypair item from the list API.")
}

func (GpuKeypair) WireDependencies(f infer.FieldSelector, _ *GpuKeypairArgs, state *GpuKeypairState) {
	f.OutputField(&state.PrivateKey).AlwaysSecret()
}
