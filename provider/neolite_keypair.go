package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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
	list, err := c.Neolite().KeypairList(ctx)
	if err != nil {
		return infer.ReadResponse[NeoliteKeypairArgs, NeoliteKeypairState]{}, err
	}
	for _, kp := range list {
		if strconv.FormatInt(kp.KeypairID, 10) == req.ID {
			state := req.State
			state.KeypairID = kp.KeypairID
			state.Name = kp.Name
			state.PublicKey = kp.PublicKey
			// private key write-only: keep value lama dari state.
			return infer.ReadResponse[NeoliteKeypairArgs, NeoliteKeypairState]{State: state}, nil
		}
	}
	return infer.ReadResponse[NeoliteKeypairArgs, NeoliteKeypairState]{},
		fmt.Errorf("neolite keypair %s not found", req.ID)
}

func (NeoliteKeypair) Delete(
	ctx context.Context, req infer.DeleteRequest[NeoliteKeypairState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	if _, err := c.Neolite().KeypairDelete(ctx, req.State.KeypairID); err != nil {
		if client.IsNotFound(err) {
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

// neoAliasStr cari value string pake beberapa kandidat key, case-insensitive.
func neoAliasStr(v map[string]any, keys ...string) string {
	for k, x := range v {
		for _, want := range keys {
			if !strings.EqualFold(k, want) {
				continue
			}
			if s, ok := x.(string); ok {
				return s
			}
		}
	}
	return ""
}

// neoAliasInt cari value int64 pake beberapa kandidat key (number atau string numerik).
func neoAliasInt(v map[string]any, keys ...string) int64 {
	for k, x := range v {
		for _, want := range keys {
			if !strings.EqualFold(k, want) {
				continue
			}
			switch n := x.(type) {
			case float64:
				return int64(n)
			case string:
				i, err := strconv.ParseInt(n, 10, 64)
				if err == nil {
					return i
				}
			}
		}
	}
	return 0
}

// neoRawJSON marshal map ke string JSON pake RedactJSON, fallback empty string.
func neoRawJSON(v map[string]any) string {
	if v == nil {
		return ""
	}
	return string(RedactJSON(v))
}
