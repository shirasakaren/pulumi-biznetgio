package provider

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
)

type ObjectStorageObject struct{}

type ObjectStorageObjectArgs struct {
	AccountID string  `pulumi:"accountId" provider:"replaceOnChanges"`
	Bucket    string  `pulumi:"bucket" provider:"replaceOnChanges"`
	Key       string  `pulumi:"key" provider:"replaceOnChanges"`
	Source    *string `pulumi:"source,optional" provider:"replaceOnChanges"`
	Content   *string `pulumi:"content,optional" provider:"secret,replaceOnChanges"`
	Acl       *string `pulumi:"acl,optional"`
}

type ObjectStorageObjectState struct {
	ObjectStorageObjectArgs
	Raw string `pulumi:"raw" provider:"secret"`
}

func (a *ObjectStorageObjectArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Object storage instance account id. Create-only.")
	ann.Describe(&a.Bucket, "Bucket name. Create-only.")
	ann.Describe(&a.Key, "Object key inside the bucket. Create-only.")
	ann.Describe(&a.Source, "Path ke file lokal yang mau di-upload. Exactly one of `source`/`content` required.")
	ann.Describe(&a.Content, "Inline content yang mau di-upload. Exactly one of `source`/`content` required.")
	ann.Describe(&a.Acl, "S3-style canned ACL applied to the object. Defaults to empty.")
	ann.SetDefault(&a.Acl, "")
}

func (s *ObjectStorageObjectState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (ObjectStorageObject) Check(
	ctx context.Context, req infer.CheckRequest,
) (infer.CheckResponse[ObjectStorageObjectArgs], error) {
	inputs, failures, err := infer.DefaultCheck[ObjectStorageObjectArgs](ctx, req.NewInputs)
	if err != nil {
		return infer.CheckResponse[ObjectStorageObjectArgs]{}, err
	}
	if (inputs.Source == nil) == (inputs.Content == nil) {
		failures = append(failures, p.CheckFailure{
			Property: "source",
			Reason:   "exactly one of source or content must be set",
		})
	}
	return infer.CheckResponse[ObjectStorageObjectArgs]{Inputs: inputs, Failures: failures}, nil
}

func (ObjectStorageObject) Create(
	ctx context.Context, req infer.CreateRequest[ObjectStorageObjectArgs],
) (infer.CreateResponse[ObjectStorageObjectState], error) {
	resp := infer.CreateResponse[ObjectStorageObjectState]{
		Output: ObjectStorageObjectState{ObjectStorageObjectArgs: req.Inputs},
	}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs
	content, err := osObjectBytes(a.Source, a.Content)
	if err != nil {
		return infer.CreateResponse[ObjectStorageObjectState]{}, err
	}
	directory, filename := osSplitKey(a.Key)
	if _, err := c.ObjectStorage().Upload(ctx, osParseID(a.AccountID), a.Bucket, directory,
		filename, content); err != nil {
		return infer.CreateResponse[ObjectStorageObjectState]{}, err
	}
	if osStr(a.Acl) != "" {
		if _, err := c.ObjectStorage().ObjectSetACL(ctx, osParseID(a.AccountID), a.Bucket, a.Key, osStr(a.Acl)); err != nil {
			return infer.CreateResponse[ObjectStorageObjectState]{}, err
		}
	}
	resp.ID = a.AccountID + ":" + a.Bucket + ":" + a.Key
	return resp, nil
}

func (ObjectStorageObject) Update(
	ctx context.Context, req infer.UpdateRequest[ObjectStorageObjectArgs, ObjectStorageObjectState],
) (infer.UpdateResponse[ObjectStorageObjectState], error) {
	resp := infer.UpdateResponse[ObjectStorageObjectState]{
		Output: ObjectStorageObjectState{ObjectStorageObjectArgs: req.Inputs},
	}
	if req.DryRun {
		return resp, nil
	}
	a := req.Inputs
	if osStr(a.Acl) == osStr(req.State.Acl) {
		return resp, nil
	}
	c := GetClient(ctx)
	if _, err := c.ObjectStorage().ObjectSetACL(ctx, osParseID(a.AccountID), a.Bucket, a.Key, osStr(a.Acl)); err != nil {
		return infer.UpdateResponse[ObjectStorageObjectState]{}, err
	}
	return resp, nil
}

func (ObjectStorageObject) Read(
	ctx context.Context, req infer.ReadRequest[ObjectStorageObjectArgs, ObjectStorageObjectState],
) (infer.ReadResponse[ObjectStorageObjectArgs, ObjectStorageObjectState], error) {
	resp := infer.ReadResponse[ObjectStorageObjectArgs, ObjectStorageObjectState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  ObjectStorageObjectState{ObjectStorageObjectArgs: req.Inputs},
	}
	c := GetClient(ctx)
	m, ok, err := osFindObject(ctx, c, osParseID(req.Inputs.AccountID), req.Inputs.Bucket, req.Inputs.Key)
	if err != nil {
		return resp, err
	}
	if !ok {
		return resp, fmt.Errorf("biznetgio: object storage object %s not found", req.ID)
	}
	if v, ok := osString(m, "acl"); ok {
		resp.State.Acl = &v
	}
	resp.State.Raw = string(osJSON(m))
	return resp, nil
}

func (ObjectStorageObject) Delete(
	ctx context.Context, req infer.DeleteRequest[ObjectStorageObjectState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	if _, err := c.ObjectStorage().ObjectDelete(ctx, osParseID(req.State.AccountID),
		req.State.Bucket, req.State.Key); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func osObjectBytes(source, content *string) ([]byte, error) {
	if source != nil && *source != "" {
		return os.ReadFile(*source)
	}
	return []byte(osStr(content)), nil
}

func osSplitKey(key string) (string, string) {
	dir, name := path.Split(key)
	return strings.TrimSuffix(dir, "/"), name
}

func osFindObject(
	ctx context.Context, c *client.Client, accountID int64, bucket, key string,
) (map[string]any, bool, error) {
	directory, name := osSplitKey(key)
	var items []map[string]any
	var err error
	if directory != "" {
		items, err = c.ObjectStorage().ObjectsListInDirectory(ctx, accountID, bucket, directory)
	} else {
		items, err = c.ObjectStorage().ObjectsList(ctx, accountID, bucket)
	}
	if err != nil {
		return nil, false, err
	}
	for _, it := range items {
		if on, ok := osString(it, "name", "key", "object_name"); ok && on == name {
			return it, true, nil
		}
	}
	return nil, false, nil
}
