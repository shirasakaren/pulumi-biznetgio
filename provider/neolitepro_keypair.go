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

func (NeoliteProKeypair) WireDependencies(
	f infer.FieldSelector, _ *NeoliteProKeypairArgs, state *NeoliteProKeypairState,
) {
	f.OutputField(&state.PrivateKey).AlwaysSecret()
}

func (NeoliteProKeypair) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteProKeypairArgs],
) (infer.CreateResponse[NeoliteProKeypairState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteProKeypairState]{
			ID:     "0",
			Output: NeoliteProKeypairState{NeoliteProKeypairArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	// pake raw response: field private key undocumented, aliasnya bisa beda-beda.
	raw, err := c.NeolitePro().KeypairCreateRaw(ctx, req.Inputs.Name)
	if err != nil {
		return infer.CreateResponse[NeoliteProKeypairState]{}, err
	}

	keypairID := neoAliasInt(raw, "keypair_id", "neosshkey_id", "id")
	if keypairID == 0 {
		return infer.CreateResponse[NeoliteProKeypairState]{},
			fmt.Errorf("create neolite pro keypair response tidak ada keypair_id: %s", neoRawJSON(raw))
	}
	state := NeoliteProKeypairState{NeoliteProKeypairArgs: req.Inputs}
	state.KeypairID = keypairID
	state.PublicKey = neoAliasStr(raw, "public_key", "pubkey")
	state.PrivateKey = neoAliasStr(raw, "private_key", "private", "secret_key", "pem")

	return infer.CreateResponse[NeoliteProKeypairState]{ID: strconv.FormatInt(keypairID, 10), Output: state}, nil
}
