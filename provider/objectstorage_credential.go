package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type ObjectStorageCredential struct{}

type ObjectStorageCredentialArgs struct {
	AccountID string `pulumi:"accountId" provider:"replaceOnChanges"`
	Active    *bool  `pulumi:"active,optional"`
}

type ObjectStorageCredentialState struct {
	ObjectStorageCredentialArgs
	AccessKey string `pulumi:"accessKey" provider:"secret"`
	SecretKey string `pulumi:"secretKey" provider:"secret"`
	Raw       string `pulumi:"raw" provider:"secret"`
}

func (a *ObjectStorageCredentialArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Object storage instance account id. Create-only.")
	ann.Describe(&a.Active, "Whether the credential is enabled. Set false to disable without deleting. Defaults to true.")
	ann.SetDefault(&a.Active, true)
}

func (s *ObjectStorageCredentialState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.AccessKey, "Credential access key, returned once at create.")
	ann.Describe(&s.SecretKey, "Credential secret key, shown only once at create; keeps its last state value on refresh.")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (ObjectStorageCredential) WireDependencies(
	f infer.FieldSelector, _ *ObjectStorageCredentialArgs, state *ObjectStorageCredentialState,
) {
	f.OutputField(&state.AccessKey).AlwaysSecret()
	f.OutputField(&state.SecretKey).AlwaysSecret()
}

func (ObjectStorageCredential) Create(
	ctx context.Context, req infer.CreateRequest[ObjectStorageCredentialArgs],
) (infer.CreateResponse[ObjectStorageCredentialState], error) {
	resp := infer.CreateResponse[ObjectStorageCredentialState]{
		Output: ObjectStorageCredentialState{ObjectStorageCredentialArgs: req.Inputs},
	}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs
	raw, err := c.ObjectStorage().CredentialCreate(ctx, osParseID(a.AccountID))
	if err != nil {
		return infer.CreateResponse[ObjectStorageCredentialState]{}, err
	}
	accessKey, ok := osString(raw, "accessKey", "access_key", "accesskey")
	if !ok {
		return infer.CreateResponse[ObjectStorageCredentialState]{},
			fmt.Errorf("biznetgio: credential create response missing access key: %s", osJSON(raw))
	}
	resp.ID = a.AccountID + ":" + osHashKey(accessKey)
	resp.Output = ObjectStorageCredentialState{
		ObjectStorageCredentialArgs: a,
		AccessKey:                   accessKey,
		SecretKey:                   osStringDefault(raw, "secretKey", "secret_key", "secretkey"),
		Raw:                         string(osJSON(raw)),
	}
	return resp, nil
}

func (ObjectStorageCredential) Update(
	ctx context.Context, req infer.UpdateRequest[ObjectStorageCredentialArgs, ObjectStorageCredentialState],
) (infer.UpdateResponse[ObjectStorageCredentialState], error) {
	resp := infer.UpdateResponse[ObjectStorageCredentialState]{Output: req.State}
	if req.DryRun {
		return resp, nil
	}
	a := req.Inputs
	if a.Active != nil && req.State.Active != nil && *a.Active != *req.State.Active {
		c := GetClient(ctx)
		accountID, keyForm, ok := osIDParts(req.ID)
		if !ok {
			return infer.UpdateResponse[ObjectStorageCredentialState]{},
				fmt.Errorf("biznetgio: invalid credential id %q", req.ID)
		}
		accessKey, ok, err := osResolveAccessKey(ctx, c, accountID, keyForm, req.State.AccessKey)
		if err != nil {
			return infer.UpdateResponse[ObjectStorageCredentialState]{}, err
		}
		if !ok {
			return infer.UpdateResponse[ObjectStorageCredentialState]{},
				fmt.Errorf("biznetgio: credential access key tidak ada di state, tidak bisa update")
		}
		if _, err := c.ObjectStorage().CredentialUpdateStatus(ctx, accountID, accessKey, *a.Active); err != nil {
			return infer.UpdateResponse[ObjectStorageCredentialState]{}, err
		}
	}
	st := req.State
	st.ObjectStorageCredentialArgs = a
	resp.Output = st
	return resp, nil
}

func (ObjectStorageCredential) Read(
	ctx context.Context, req infer.ReadRequest[ObjectStorageCredentialArgs, ObjectStorageCredentialState],
) (infer.ReadResponse[ObjectStorageCredentialArgs, ObjectStorageCredentialState], error) {
	resp := infer.ReadResponse[ObjectStorageCredentialArgs, ObjectStorageCredentialState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  ObjectStorageCredentialState{ObjectStorageCredentialArgs: req.Inputs},
	}
	c := GetClient(ctx)
