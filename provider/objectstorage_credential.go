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
