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

