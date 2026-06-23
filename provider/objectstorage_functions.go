package provider

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ObjectStorageInstances struct{}

type ObjectStorageInstancesArgs struct {
	Status *string `pulumi:"status,optional"`
}

type ObjectStorageInstancesResult struct {
	Instances []ObjectStorageInstancesInstance `pulumi:"instances"`
}

type ObjectStorageInstancesInstance struct {
	ID        string `pulumi:"id"`
	Label     string `pulumi:"label"`
	Status    string `pulumi:"status"`
	ProductID int64  `pulumi:"productId"`
	Quota     int64  `pulumi:"quota"`
	Raw       string `pulumi:"raw" provider:"secret"`
}

func (a *ObjectStorageInstancesArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Status, "Filter by status (Active, Pending, Suspended, Terminated).")
}

func (r *ObjectStorageInstancesResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Instances, "List of object storage instances.")
}

func (ObjectStorageInstances) Invoke(
	ctx context.Context, req infer.FunctionRequest[ObjectStorageInstancesArgs],
) (infer.FunctionResponse[ObjectStorageInstancesResult], error) {
	c := GetClient(ctx)
	items, err := c.ObjectStorage().AccountsList(ctx, osStr(req.Input.Status))
	if err != nil {
		return infer.FunctionResponse[ObjectStorageInstancesResult]{}, err
	}
	out := make([]ObjectStorageInstancesInstance, 0, len(items))
	for _, it := range items {
		i := ObjectStorageInstancesInstance{
			Label:  osStringDefault(it, "label", "name"),
			Status: osStringDefault(it, "status", "state"),
			Raw:    string(osJSON(it)),
		}
		if v, ok := osString(it, "account_id", "id"); ok {
			i.ID = v
		}
		if v, ok := osInt(it, "product_id"); ok {
			i.ProductID = v
		}
		if v, ok := osInt(it, "quota"); ok {
			i.Quota = v
		}
		out = append(out, i)
	}
	return infer.FunctionResponse[ObjectStorageInstancesResult]{Output: ObjectStorageInstancesResult{Instances: out}}, nil
}

type ObjectStorageBuckets struct{}

type ObjectStorageBucketsArgs struct {
	AccountID string `pulumi:"accountId"`
}

type ObjectStorageBucketsResult struct {
	Buckets []ObjectStorageBucketsBucket `pulumi:"buckets"`
}

type ObjectStorageBucketsBucket struct {
	Name string `pulumi:"name"`
	ACL  string `pulumi:"acl"`
	Raw  string `pulumi:"raw" provider:"secret"`
}

func (a *ObjectStorageBucketsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Object storage instance account id.")
}

func (r *ObjectStorageBucketsResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Buckets, "List of buckets inside the instance.")
}

func (ObjectStorageBuckets) Invoke(
	ctx context.Context, req infer.FunctionRequest[ObjectStorageBucketsArgs],
) (infer.FunctionResponse[ObjectStorageBucketsResult], error) {
	c := GetClient(ctx)
	items, err := c.ObjectStorage().BucketsList(ctx, osParseID(req.Input.AccountID))
	if err != nil {
		return infer.FunctionResponse[ObjectStorageBucketsResult]{}, err
	}
	out := make([]ObjectStorageBucketsBucket, 0, len(items))
	for _, it := range items {
		out = append(out, ObjectStorageBucketsBucket{
			Name: osStringDefault(it, "name", "bucket_name"),
			ACL:  osStringDefault(it, "acl"),
			Raw:  string(osJSON(it)),
		})
	}
	return infer.FunctionResponse[ObjectStorageBucketsResult]{Output: ObjectStorageBucketsResult{Buckets: out}}, nil
}

type ObjectStorageCredentials struct{}

type ObjectStorageCredentialsArgs struct {
	AccountID string `pulumi:"accountId"`
}

type ObjectStorageCredentialsResult struct {
	Credentials []ObjectStorageCredentialsCredential `pulumi:"credentials"`
}

type ObjectStorageCredentialsCredential struct {
	AccessKey string `pulumi:"accessKey" provider:"secret"`
	Active    bool   `pulumi:"active"`
}

func (a *ObjectStorageCredentialsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Object storage instance account id.")
}

func (r *ObjectStorageCredentialsResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Credentials, "List of credentials (access key + active state only; "+
		"secret keys tidak pernah dikeluarkan).")
}

func (ObjectStorageCredentials) Invoke(
	ctx context.Context, req infer.FunctionRequest[ObjectStorageCredentialsArgs],
) (infer.FunctionResponse[ObjectStorageCredentialsResult], error) {
	c := GetClient(ctx)
	items, err := c.ObjectStorage().CredentialsList(ctx, osParseID(req.Input.AccountID))
	if err != nil {
		return infer.FunctionResponse[ObjectStorageCredentialsResult]{}, err
	}
	out := make([]ObjectStorageCredentialsCredential, 0, len(items))
	for _, it := range items {
		cred := ObjectStorageCredentialsCredential{
			AccessKey: osStringDefault(it, "accessKey", "access_key", "accesskey"),
		}
		if v, ok := osBool(it, "active"); ok {
			cred.Active = v
		}
		out = append(out, cred)
	}
	return infer.FunctionResponse[ObjectStorageCredentialsResult]{
		Output: ObjectStorageCredentialsResult{Credentials: out},
	}, nil
}

var (
	_ infer.Fn[ObjectStorageInstancesArgs, ObjectStorageInstancesResult]     = ObjectStorageInstances{}
	_ infer.Fn[ObjectStorageBucketsArgs, ObjectStorageBucketsResult]         = ObjectStorageBuckets{}
	_ infer.Fn[ObjectStorageCredentialsArgs, ObjectStorageCredentialsResult] = ObjectStorageCredentials{}
)
