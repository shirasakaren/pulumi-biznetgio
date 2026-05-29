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
// wip 281
