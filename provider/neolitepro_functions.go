package provider

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// ---------- getNeoliteProProducts ----------

type NeoliteProProducts struct{}

type NeoliteProProductsArgs struct{}

type NeoliteProProductsResult struct {
	Products []NeolitePlanItem `pulumi:"products"`
}

func (NeoliteProProducts) Invoke(
	ctx context.Context, _ infer.FunctionRequest[NeoliteProProductsArgs],
) (infer.FunctionResponse[NeoliteProProductsResult], error) {
	c := GetClient(ctx)
	plans, err := c.NeolitePro().ProductList(ctx)
	if err != nil {
		return infer.FunctionResponse[NeoliteProProductsResult]{}, err
	}

	result := NeoliteProProductsResult{}
	for _, p := range plans {
		item := NeolitePlanItem{
			ProductID:    p.ProductID,
			Name:         p.Name,
			Description:  p.Description,
			CategoryID:   p.CategoryID,
			CategoryName: p.CategoryName,
			Options: NeolitePlanOptions{
				Type:           p.Options.Type,
				Cores:          p.Options.Cores,
