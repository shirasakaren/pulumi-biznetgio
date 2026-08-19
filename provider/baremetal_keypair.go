package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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
		"Only appears once; preserved across refreshes.")
	ann.Describe(&s.Raw, "Raw JSON of the keypair item from the list API.")
}

func (BaremetalKeypair) WireDependencies(
	f infer.FieldSelector, _ *BaremetalKeypairArgs, state *BaremetalKeypairState,
) {
	f.OutputField(&state.PrivateKey).AlwaysSecret()
}

func (BaremetalKeypair) Create(
	ctx context.Context, req infer.CreateRequest[BaremetalKeypairArgs],
) (infer.CreateResponse[BaremetalKeypairState], error) {
	resp := infer.CreateResponse[BaremetalKeypairState]{Output: keypairStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs

	var raw map[string]any
	var err error
	if bmtStr(a.PublicKey) == "" {
		raw, err = c.Baremetal().KeypairCreate(ctx, a.Name)
	} else {
		raw, err = c.Baremetal().KeypairImport(ctx, client.KeypairImportPayload{Name: a.Name, PublicKey: *a.PublicKey})
	}
	if err != nil {
		return infer.CreateResponse[BaremetalKeypairState]{}, err
	}
	keypairID, ok := bmtInt64(raw, "keypair_id", "id")
	if !ok {
		return infer.CreateResponse[BaremetalKeypairState]{},
			fmt.Errorf("biznetgio: keypair create response missing id: %s", bmtJSON(raw))
	}
	resp.ID = strconv.FormatInt(keypairID, 10)
	resp.Output = keypairStateFromMap(ctx, a, raw)
	return resp, nil
}

func (BaremetalKeypair) Read(
	ctx context.Context, req infer.ReadRequest[BaremetalKeypairArgs, BaremetalKeypairState],
) (infer.ReadResponse[BaremetalKeypairArgs, BaremetalKeypairState], error) {
	resp := infer.ReadResponse[BaremetalKeypairArgs, BaremetalKeypairState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  keypairStateFromMap(ctx, req.Inputs, nil),
	}
	c := GetClient(ctx)
	keypairID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return resp, fmt.Errorf("biznetgio: invalid keypair id %q: %s", req.ID, err)
	}
	items, err := c.Baremetal().KeypairList(ctx)
	if err != nil {
		return resp, err
	}
	for _, it := range items {
		if id, ok := bmtInt64(it, "keypair_id", "id"); ok && id == keypairID {
			resp.State = keypairStateFromMap(ctx, req.Inputs, it)
			// list gak bawa private key - keep value lama
			resp.State.PrivateKey = req.State.PrivateKey
			return resp, nil
		}
	}
	return resp, fmt.Errorf("biznetgio: keypair %s not found", req.ID)
}

func (BaremetalKeypair) Delete(
	ctx context.Context, req infer.DeleteRequest[BaremetalKeypairState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	keypairID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("biznetgio: invalid keypair id %q: %s", req.ID, err)
	}
	if _, err := c.Baremetal().KeypairDelete(ctx, keypairID); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func keypairStateFromMap(_ context.Context, args BaremetalKeypairArgs, m map[string]any) BaremetalKeypairState {
	st := BaremetalKeypairState{BaremetalKeypairArgs: args}
	if m == nil {
		return st
	}
	if v, ok := bmtInt64(m, "keypair_id", "id"); ok {
		st.KeypairID = v
	}
	if v, ok := bmtString(m, "name"); ok {
		st.Name = v
	}
	if v, ok := bmtString(m, "public_key", "publickey"); ok {
		st.PublicKey = &v
	}
	// private key cuma di response create - alias defensif, jangan nebak nama
	if v, ok := bmtString(m, "private_key", "private", "secret_key", "pem"); ok {
		st.PrivateKey = &v
	}
	st.Raw = string(bmtJSON(m))
	return st
}
