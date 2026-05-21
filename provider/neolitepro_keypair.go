package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type NeoliteProKeypair struct{}

type NeoliteProKeypairArgs struct {
	Name string `pulumi:"name"`
}

type NeoliteProKeypairState struct {
	NeoliteProKeypairArgs
	KeypairID  int64  `pulumi:"keypairId"`
	PublicKey  string `pulumi:"publicKey"`
	PrivateKey string `pulumi:"privateKey" provider:"secret"`
}

func (a *NeoliteProKeypairArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Nama keypair.")
}

func (s *NeoliteProKeypairState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.KeypairID, "Id keypair di BiznetGIO.")
	ann.Describe(&s.PublicKey, "Public key yang di-generate.")
	ann.Describe(&s.PrivateKey, "Private key (sensitive). Write-only: cuma ada di response create, ga bisa di-refetch.")
}

