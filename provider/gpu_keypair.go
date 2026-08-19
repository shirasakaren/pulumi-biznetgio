package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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

func (GpuKeypair) Create(
	ctx context.Context, req infer.CreateRequest[GpuKeypairArgs],
) (infer.CreateResponse[GpuKeypairState], error) {
	resp := infer.CreateResponse[GpuKeypairState]{Output: GpuKeypairState{GpuKeypairArgs: req.Inputs}}
	if req.DryRun {
		return resp, nil
	}
	m, err := GetClient(ctx).GPU().KeypairCreate(ctx, req.Inputs.Name)
	if err != nil {
		return infer.CreateResponse[GpuKeypairState]{}, err
	}
	keypairID, ok := gpuInt64(m, "keypair_id", "id")
	if !ok {
		return infer.CreateResponse[GpuKeypairState]{},
			fmt.Errorf("biznetgio: keypair create response missing keypair_id: %s", gpuJSON(m))
	}
	resp.ID = strconv.FormatInt(keypairID, 10)
	resp.Output = gpuKeypairStateFromMap(req.Inputs, m)
	return resp, nil
}

func (GpuKeypair) Read(
	ctx context.Context, req infer.ReadRequest[GpuKeypairArgs, GpuKeypairState],
) (infer.ReadResponse[GpuKeypairArgs, GpuKeypairState], error) {
	resp := infer.ReadResponse[GpuKeypairArgs, GpuKeypairState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  GpuKeypairState{GpuKeypairArgs: req.Inputs},
	}
	c := GetClient(ctx)
	keypairID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return resp, fmt.Errorf("biznetgio: invalid keypair id %q: %s", req.ID, err)
	}
	list, err := c.GPU().KeypairList(ctx)
	if err != nil {
		return resp, err
	}
	for _, m := range list {
		id, ok := gpuInt64(m, "keypair_id", "id")
		if ok && id == keypairID {
			resp.State = gpuKeypairStateFromMap(req.Inputs, m)
			return resp, nil
		}
	}
	return resp, fmt.Errorf("biznetgio: keypair %s not found", req.ID)
}

func (GpuKeypair) Delete(ctx context.Context, req infer.DeleteRequest[GpuKeypairState]) (infer.DeleteResponse, error) {
	keypairID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("biznetgio: invalid keypair id %q: %s", req.ID, err)
	}
	if _, err := GetClient(ctx).GPU().KeypairDelete(ctx, keypairID); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func gpuKeypairStateFromMap(args GpuKeypairArgs, m map[string]any) GpuKeypairState {
	st := GpuKeypairState{GpuKeypairArgs: args}
	if v, ok := gpuString(m, "name"); ok {
		st.Name = v
	}
	if v, ok := gpuString(m, "public_key", "pubkey"); ok {
		st.PublicKey = &v
	}
	if v, ok := gpuString(m, "private_key", "private", "secret_key", "pem"); ok {
		st.PrivateKey = &v
	}
	st.Raw = string(gpuJSON(m))
	return st
}
