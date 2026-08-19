package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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

func (NeoliteProKeypair) Read(
	ctx context.Context, req infer.ReadRequest[NeoliteProKeypairArgs, NeoliteProKeypairState],
) (infer.ReadResponse[NeoliteProKeypairArgs, NeoliteProKeypairState], error) {
	c := GetClient(ctx)
	list, err := c.NeolitePro().KeypairList(ctx)
	if err != nil {
		return infer.ReadResponse[NeoliteProKeypairArgs, NeoliteProKeypairState]{}, err
	}
	for _, kp := range list {
		if strconv.FormatInt(kp.KeypairID, 10) == req.ID {
			state := req.State
			state.KeypairID = kp.KeypairID
			state.Name = kp.Name
			state.PublicKey = kp.PublicKey
			// private key write-only: keep value lama dari state.
			return infer.ReadResponse[NeoliteProKeypairArgs, NeoliteProKeypairState]{State: state}, nil
		}
	}
	return infer.ReadResponse[NeoliteProKeypairArgs, NeoliteProKeypairState]{},
		fmt.Errorf("neolite pro keypair %s not found", req.ID)
}

func (NeoliteProKeypair) Delete(
	ctx context.Context, req infer.DeleteRequest[NeoliteProKeypairState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	if _, err := c.NeolitePro().KeypairDelete(ctx, req.State.KeypairID); err != nil {
		if client.IsNotFound(err) {
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
