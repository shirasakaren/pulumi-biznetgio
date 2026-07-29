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
