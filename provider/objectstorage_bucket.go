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
