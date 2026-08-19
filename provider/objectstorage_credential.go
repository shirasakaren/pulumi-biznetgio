package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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
				fmt.Errorf("biznetgio: credential access key missing from state, can't update")
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
	accountID, keyForm, ok := osIDParts(req.ID)
	if !ok {
		return resp, fmt.Errorf("biznetgio: invalid credential id %q", req.ID)
	}
	m, ok, err := osFindCredential(ctx, c, accountID, keyForm)
	if err != nil {
		return resp, err
	}
	if !ok {
		return resp, fmt.Errorf("biznetgio: object storage credential %s not found", req.ID)
	}
	accessKey, _ := osString(m, "accessKey", "access_key", "accesskey")
	resp.State.AccessKey = accessKey
	if v, ok := osString(m, "secretKey", "secret_key", "secretkey"); ok {
		resp.State.SecretKey = v
	}
	if v, ok := osBool(m, "active"); ok {
		resp.State.Active = &v
	}
	resp.State.Raw = string(osJSON(m))
	return resp, nil
}

func (ObjectStorageCredential) Delete(
	ctx context.Context, req infer.DeleteRequest[ObjectStorageCredentialState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	accountID, keyForm, ok := osIDParts(req.ID)
	if !ok {
		return infer.DeleteResponse{}, fmt.Errorf("biznetgio: invalid credential id %q", req.ID)
	}
	accessKey, ok, err := osResolveAccessKey(ctx, c, accountID, keyForm, req.State.AccessKey)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if !ok {
		return infer.DeleteResponse{}, fmt.Errorf("biznetgio: credential access key tidak ada di state")
	}
	if _, err := c.ObjectStorage().CredentialDelete(ctx, accountID, accessKey); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func osFindCredential(
	ctx context.Context, c *client.Client, accountID int64, keyForm string,
) (map[string]any, bool, error) {
	items, err := c.ObjectStorage().CredentialsList(ctx, accountID)
	if err != nil {
		return nil, false, err
	}
	hashMode := osIsHex16(keyForm)
	for _, it := range items {
		ak, ok := osString(it, "accessKey", "access_key", "accesskey")
		if !ok {
			continue
		}
		if (hashMode && osHashKey(ak) == keyForm) || (!hashMode && ak == keyForm) {
			return it, true, nil
		}
	}
	return nil, false, nil
}

// osResolveAccessKey returns the plaintext access key from keyForm (hash or literal).
// Hash mode needs the credentials list - falls back if not found (old state).
func osResolveAccessKey(
	ctx context.Context, c *client.Client, accountID int64, keyForm, fallback string,
) (string, bool, error) {
	if !osIsHex16(keyForm) {
		return keyForm, keyForm != "", nil
	}
	items, err := c.ObjectStorage().CredentialsList(ctx, accountID)
	if err != nil {
		return "", false, err
	}
	for _, it := range items {
		if ak, ok := osString(it, "accessKey", "access_key", "accesskey"); ok && osHashKey(ak) == keyForm {
			return ak, true, nil
		}
	}
	if fallback != "" {
		return fallback, true, nil
	}
	return "", false, nil
}

func osHashKey(accessKey string) string {
	sum := sha256.Sum256([]byte(accessKey))
	return hex.EncodeToString(sum[:8])
}

func osIsHex16(s string) bool {
	if len(s) != 16 {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}
	return true
}

func osIDParts(id string) (int64, string, bool) {
	acc, k, ok := strings.Cut(id, ":")
	if !ok || k == "" {
		return 0, "", false
	}
	n, err := strconv.ParseInt(acc, 10, 64)
	if err != nil || n == 0 {
		return 0, "", false
	}
	return n, k, true
}
