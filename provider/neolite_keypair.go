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
	f.OutputField(&state.PrivateKey).AlwaysSecret()
}

func (NeoliteKeypair) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteKeypairArgs],
) (infer.CreateResponse[NeoliteKeypairState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteKeypairState]{
			ID:     "0",
			Output: NeoliteKeypairState{NeoliteKeypairArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	// pake raw response: field private key undocumented, aliasnya bisa beda-beda.
	raw, err := c.Neolite().KeypairCreateRaw(ctx, req.Inputs.Name)
	if err != nil {
		return infer.CreateResponse[NeoliteKeypairState]{}, err
	}

	keypairID := neoAliasInt(raw, "keypair_id", "neosshkey_id", "id")
	if keypairID == 0 {
		return infer.CreateResponse[NeoliteKeypairState]{},
			fmt.Errorf("create neolite keypair response tidak ada keypair_id: %s", neoRawJSON(raw))
	}
	state := NeoliteKeypairState{NeoliteKeypairArgs: req.Inputs}
	state.KeypairID = keypairID
	state.PublicKey = neoAliasStr(raw, "public_key", "pubkey")
	state.PrivateKey = neoAliasStr(raw, "private_key", "private", "secret_key", "pem")

	return infer.CreateResponse[NeoliteKeypairState]{ID: strconv.FormatInt(keypairID, 10), Output: state}, nil
}

func (NeoliteKeypair) Read(
	ctx context.Context, req infer.ReadRequest[NeoliteKeypairArgs, NeoliteKeypairState],
) (infer.ReadResponse[NeoliteKeypairArgs, NeoliteKeypairState], error) {
	c := GetClient(ctx)
