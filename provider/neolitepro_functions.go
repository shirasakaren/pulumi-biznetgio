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
				Memory:         p.Options.Memory,
				AllowDowngrade: p.Options.AllowDowngrade,
			},
		}
		for _, b := range p.Billing {
			bm := NeolitePlanBilling{Label: b.Label, Cycle: b.Cycle, Price: b.Price}
			for _, c := range b.Components {
				cm := NeolitePlanComponent{Label: c.Label, Field: c.Field}
				for _, pr := range c.Prices {
					cm.Prices = append(cm.Prices, NeolitePlanPrice{QtyMin: pr.QtyMin, QtyMax: pr.QtyMax, Price: pr.Price})
				}
				bm.Components = append(bm.Components, cm)
			}
			item.Billing = append(item.Billing, bm)
		}
		result.Products = append(result.Products, item)
	}
	return infer.FunctionResponse[NeoliteProProductsResult]{Output: result}, nil
}

// ---------- getNeoliteProOsList ----------

type NeoliteProOsList struct{}

type NeoliteProOsListArgs struct {
	ProductID int64 `pulumi:"productId"`
}

type NeoliteProOsListResult struct {
	Oss []NeoliteOsItem `pulumi:"oss"`
}

func (a *NeoliteProOsListArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id NEO Lite Pro.")
}

func (NeoliteProOsList) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteProOsListArgs],
) (infer.FunctionResponse[NeoliteProOsListResult], error) {
	c := GetClient(ctx)
	oss, err := c.NeolitePro().ProductOSList(ctx, req.Input.ProductID)
	if err != nil {
		return infer.FunctionResponse[NeoliteProOsListResult]{}, err
	}

	result := NeoliteProOsListResult{}
	for _, os := range oss {
		result.Oss = append(result.Oss, NeoliteOsItem{
			VMID:   os.VMID,
			Node:   os.Node,
			Name:   os.Name,
			MaxMem: os.MaxMem,
			MaxCPU: os.MaxCPU,
		})
	}
	return infer.FunctionResponse[NeoliteProOsListResult]{Output: result}, nil
}

// ---------- getNeoliteProChangePackageOptions ----------

type NeoliteProChangePackageOptions struct{}

type NeoliteProChangePackageOptionsArgs struct {
	AccountID int64 `pulumi:"accountId"`
}

type NeoliteProChangePackageOptionsResult struct {
	Raw string `pulumi:"raw" provider:"secret"`
}

func (a *NeoliteProChangePackageOptionsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Account id VM NEO Lite Pro.")
}

func (NeoliteProChangePackageOptions) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteProChangePackageOptionsArgs],
