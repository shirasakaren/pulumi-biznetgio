package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type NeoliteKeypair struct{}

type NeoliteKeypairArgs struct {
	Name string `pulumi:"name"`
}

type NeoliteKeypairState struct {
	NeoliteKeypairArgs
	KeypairID  int64  `pulumi:"keypairId"`
	PublicKey  string `pulumi:"publicKey"`
	PrivateKey string `pulumi:"privateKey" provider:"secret"`
}

func (a *NeoliteKeypairArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Name, "Nama keypair.")
}

func (s *NeoliteKeypairState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.KeypairID, "Id keypair di BiznetGIO.")
	ann.Describe(&s.PublicKey, "Public key yang di-generate.")
	ann.Describe(&s.PrivateKey, "Private key (sensitive). Write-only: cuma ada di response create, ga bisa di-refetch.")
}

func (NeoliteKeypair) WireDependencies(f infer.FieldSelector, _ *NeoliteKeypairArgs, state *NeoliteKeypairState) {
