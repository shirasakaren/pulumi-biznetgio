package provider

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
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

