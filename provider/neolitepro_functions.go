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
	ann.Describe(&a.ProductID, "NEO Lite Pro product id.")
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
	ann.Describe(&a.AccountID, "Account id of the NEO Lite Pro VM.")
}

func (NeoliteProChangePackageOptions) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteProChangePackageOptionsArgs],
) (infer.FunctionResponse[NeoliteProChangePackageOptionsResult], error) {
	c := GetClient(ctx)
	out, err := c.NeolitePro().ChangePackageOptions(ctx, req.Input.AccountID)
	if err != nil {
		return infer.FunctionResponse[NeoliteProChangePackageOptionsResult]{}, err
	}
	return infer.FunctionResponse[NeoliteProChangePackageOptionsResult]{
		Output: NeoliteProChangePackageOptionsResult{Raw: neoRawJSON(out)},
	}, nil
}

// ---------- getNeoliteProStorageUpgradeOptions ----------

type NeoliteProStorageUpgradeOptions struct{}

type NeoliteProStorageUpgradeOptionsArgs struct {
	AccountID int64 `pulumi:"accountId"`
}

type NeoliteProStorageUpgradeOptionsResult struct {
	Raw string `pulumi:"raw" provider:"secret"`
}

func (a *NeoliteProStorageUpgradeOptionsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Account id of the NEO Lite Pro VM.")
}

func (NeoliteProStorageUpgradeOptions) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteProStorageUpgradeOptionsArgs],
) (infer.FunctionResponse[NeoliteProStorageUpgradeOptionsResult], error) {
	c := GetClient(ctx)
	out, err := c.NeolitePro().StorageOptions(ctx, req.Input.AccountID)
	if err != nil {
		return infer.FunctionResponse[NeoliteProStorageUpgradeOptionsResult]{}, err
	}
	return infer.FunctionResponse[NeoliteProStorageUpgradeOptionsResult]{
		Output: NeoliteProStorageUpgradeOptionsResult{Raw: neoRawJSON(out)},
	}, nil
}

// ---------- getNeoliteProIPAvailability ----------

type NeoliteProIPAvailability struct{}

type NeoliteProIPAvailabilityArgs struct {
	ProductID int64 `pulumi:"productId"`
}

type NeoliteProIPAvailabilityResult struct {
	Available bool `pulumi:"available"`
}

func (a *NeoliteProIPAvailabilityArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "NEO Lite Pro product id.")
}

func (NeoliteProIPAvailability) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteProIPAvailabilityArgs],
) (infer.FunctionResponse[NeoliteProIPAvailabilityResult], error) {
	c := GetClient(ctx)
	out, err := c.NeolitePro().ProductIPAvailability(ctx, req.Input.ProductID)
	if err != nil {
		return infer.FunctionResponse[NeoliteProIPAvailabilityResult]{}, err
	}
	return infer.FunctionResponse[NeoliteProIPAvailabilityResult]{
		Output: NeoliteProIPAvailabilityResult{Available: out.Available},
	}, nil
}

var (
	_ infer.Fn[NeoliteProProductsArgs, NeoliteProProductsResult]                           = NeoliteProProducts{}
	_ infer.Fn[NeoliteProOsListArgs, NeoliteProOsListResult]                               = NeoliteProOsList{}
	_ infer.Fn[NeoliteProChangePackageOptionsArgs, NeoliteProChangePackageOptionsResult]   = NeoliteProChangePackageOptions{}  //nolint:lll // gofmt merge balik
	_ infer.Fn[NeoliteProStorageUpgradeOptionsArgs, NeoliteProStorageUpgradeOptionsResult] = NeoliteProStorageUpgradeOptions{} //nolint:lll // gofmt merge balik
	_ infer.Fn[NeoliteProIPAvailabilityArgs, NeoliteProIPAvailabilityResult]               = NeoliteProIPAvailability{}
)
