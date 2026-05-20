package provider

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type ObjectStorageBucket struct{}

type ObjectStorageBucketArgs struct {
	AccountID string  `pulumi:"accountId" provider:"replaceOnChanges"`
	Name      string  `pulumi:"name" provider:"replaceOnChanges"`
	Acl       *string `pulumi:"acl,optional"`
}

type ObjectStorageBucketState struct {
	ObjectStorageBucketArgs
	Raw string `pulumi:"raw" provider:"secret"`
}

func (a *ObjectStorageBucketArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Object storage instance account id. Create-only.")
	ann.Describe(&a.Name, "Bucket name, minimum 3 chars, `[a-zA-Z0-9-_]`. Immutable, changing it replaces the bucket.")
	ann.Describe(&a.Acl, "S3-style canned ACL: empty (default private), private, public-read, "+
		"public-read-write, authenticated-read, log-delivery-write. Defaults to empty.")
	ann.SetDefault(&a.Acl, "")
}

func (s *ObjectStorageBucketState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (ObjectStorageBucket) Create(
	ctx context.Context, req infer.CreateRequest[ObjectStorageBucketArgs],
) (infer.CreateResponse[ObjectStorageBucketState], error) {
	resp := infer.CreateResponse[ObjectStorageBucketState]{
		Output: ObjectStorageBucketState{ObjectStorageBucketArgs: req.Inputs},
	}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs
	if _, err := c.ObjectStorage().BucketCreate(ctx, osParseID(a.AccountID), client.NOSCreateBucketPayload{
		Name: a.Name,
		ACL:  osStr(a.Acl),
	}); err != nil {
		return infer.CreateResponse[ObjectStorageBucketState]{}, err
	}
	resp.ID = a.AccountID + ":" + a.Name
	return resp, nil
}

func (ObjectStorageBucket) Update(
	ctx context.Context, req infer.UpdateRequest[ObjectStorageBucketArgs, ObjectStorageBucketState],
) (infer.UpdateResponse[ObjectStorageBucketState], error) {
	resp := infer.UpdateResponse[ObjectStorageBucketState]{
		Output: ObjectStorageBucketState{ObjectStorageBucketArgs: req.Inputs},
	}
	if req.DryRun {
		return resp, nil
	}
	a := req.Inputs
	if osStr(a.Acl) == osStr(req.State.Acl) {
		return resp, nil
	}
	c := GetClient(ctx)
	if _, err := c.ObjectStorage().BucketSetACL(ctx, osParseID(a.AccountID), a.Name, osStr(a.Acl)); err != nil {
		return infer.UpdateResponse[ObjectStorageBucketState]{}, err
	}
	return resp, nil
}

func (ObjectStorageBucket) Read(
	ctx context.Context, req infer.ReadRequest[ObjectStorageBucketArgs, ObjectStorageBucketState],
) (infer.ReadResponse[ObjectStorageBucketArgs, ObjectStorageBucketState], error) {
	resp := infer.ReadResponse[ObjectStorageBucketArgs, ObjectStorageBucketState]{
		ID:     req.ID,
