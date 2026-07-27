package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type BaremetalKeypair struct{}

type BaremetalKeypairArgs struct {
	Name      string  `pulumi:"name"`
	PublicKey *string `pulumi:"publicKey,optional"`
}

type BaremetalKeypairState struct {
	BaremetalKeypairArgs
	KeypairID  int64   `pulumi:"keypairId"`
	PrivateKey *string `pulumi:"privateKey" provider:"secret"`
	Raw        string  `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalKeypairArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Name of the keypair.")
	ann.Describe(&a.PublicKey, "When set, the keypair is imported via `POST /baremetals/keypairs/import`. "+
		"When empty, the keypair is generated server-side. Create-only, changing it replaces the keypair.")
}

func (s *BaremetalKeypairState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.KeypairID, "Keypair id in BiznetGIO.")
	ann.Describe(&s.PrivateKey, "Private key from the creation response, when the API returns one. "+
// wip 808
