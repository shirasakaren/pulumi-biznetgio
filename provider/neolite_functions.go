package provider

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// ---------- getProducts ----------

type NeoliteProducts struct{}

type NeoliteProductsArgs struct{}

type NeoliteProductsResult struct {
	Products []NeolitePlanItem `pulumi:"products"`
}

type NeolitePlanItem struct {
	ProductID    int64                `pulumi:"productId"`
	Name         string               `pulumi:"name"`
	Description  string               `pulumi:"description"`
	CategoryID   int64                `pulumi:"categoryId"`
	CategoryName string               `pulumi:"categoryName"`
	Options      NeolitePlanOptions   `pulumi:"options"`
	Billing      []NeolitePlanBilling `pulumi:"billing"`
}

type NeolitePlanOptions struct {
	Type           string `pulumi:"type"`
	Cores          int64  `pulumi:"cores"`
	Memory         int64  `pulumi:"memory"`
	AllowDowngrade int64  `pulumi:"allowDowngrade"`
}

type NeolitePlanBilling struct {
	Label      string                 `pulumi:"label"`
	Cycle      string                 `pulumi:"cycle"`
	Price      int64                  `pulumi:"price"`
	Components []NeolitePlanComponent `pulumi:"components"`
}

type NeolitePlanComponent struct {
	Label  string             `pulumi:"label"`
	Field  string             `pulumi:"field"`
	Prices []NeolitePlanPrice `pulumi:"prices"`
}

type NeolitePlanPrice struct {
	QtyMin int64 `pulumi:"qtyMin"`
	QtyMax int64 `pulumi:"qtyMax"`
	Price  int64 `pulumi:"price"`
}

func (NeoliteProducts) Invoke(
	ctx context.Context, _ infer.FunctionRequest[NeoliteProductsArgs],
) (infer.FunctionResponse[NeoliteProductsResult], error) {
	c := GetClient(ctx)
	plans, err := c.Neolite().ProductList(ctx)
	if err != nil {
		return infer.FunctionResponse[NeoliteProductsResult]{}, err
	}

	result := NeoliteProductsResult{}
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
	return infer.FunctionResponse[NeoliteProductsResult]{Output: result}, nil
}

// ---------- getOsList ----------

type NeoliteOsList struct{}

type NeoliteOsListArgs struct {
	ProductID int64 `pulumi:"productId"`
}

type NeoliteOsListResult struct {
	Oss []NeoliteOsItem `pulumi:"oss"`
}

type NeoliteOsItem struct {
	VMID   int64  `pulumi:"vmid"`
	Node   string `pulumi:"node"`
	Name   string `pulumi:"name"`
	MaxMem int64  `pulumi:"maxmem"`
	MaxCPU int64  `pulumi:"maxcpu"`
}

func (a *NeoliteOsListArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id NEO Lite.")
}

func (NeoliteOsList) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteOsListArgs],
) (infer.FunctionResponse[NeoliteOsListResult], error) {
	c := GetClient(ctx)
	oss, err := c.Neolite().ProductOSList(ctx, req.Input.ProductID)
	if err != nil {
		return infer.FunctionResponse[NeoliteOsListResult]{}, err
	}

	result := NeoliteOsListResult{}
	for _, os := range oss {
		result.Oss = append(result.Oss, NeoliteOsItem{
			VMID:   os.VMID,
			Node:   os.Node,
			Name:   os.Name,
			MaxMem: os.MaxMem,
			MaxCPU: os.MaxCPU,
		})
	}
	return infer.FunctionResponse[NeoliteOsListResult]{Output: result}, nil
}

// ---------- getChangePackageOptions ----------

type NeoliteChangePackageOptions struct{}

type NeoliteChangePackageOptionsArgs struct {
	AccountID int64 `pulumi:"accountId"`
}

type NeoliteChangePackageOptionsResult struct {
	Raw string `pulumi:"raw" provider:"secret"`
}

func (a *NeoliteChangePackageOptionsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Account id VM NEO Lite.")
}

func (NeoliteChangePackageOptions) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteChangePackageOptionsArgs],
) (infer.FunctionResponse[NeoliteChangePackageOptionsResult], error) {
	c := GetClient(ctx)
	out, err := c.Neolite().ChangePackageOptions(ctx, req.Input.AccountID)
	if err != nil {
		return infer.FunctionResponse[NeoliteChangePackageOptionsResult]{}, err
	}
	return infer.FunctionResponse[NeoliteChangePackageOptionsResult]{
		Output: NeoliteChangePackageOptionsResult{Raw: neoRawJSON(out)},
	}, nil
}

// ---------- getStorageUpgradeOptions ----------

type NeoliteStorageUpgradeOptions struct{}

type NeoliteStorageUpgradeOptionsArgs struct {
	AccountID int64 `pulumi:"accountId"`
}

type NeoliteStorageUpgradeOptionsResult struct {
	Raw string `pulumi:"raw" provider:"secret"`
}

func (a *NeoliteStorageUpgradeOptionsArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Account id VM NEO Lite.")
}

func (NeoliteStorageUpgradeOptions) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteStorageUpgradeOptionsArgs],
) (infer.FunctionResponse[NeoliteStorageUpgradeOptionsResult], error) {
	c := GetClient(ctx)
	out, err := c.Neolite().StorageOptions(ctx, req.Input.AccountID)
	if err != nil {
		return infer.FunctionResponse[NeoliteStorageUpgradeOptionsResult]{}, err
	}
	return infer.FunctionResponse[NeoliteStorageUpgradeOptionsResult]{
		Output: NeoliteStorageUpgradeOptionsResult{Raw: neoRawJSON(out)},
	}, nil
}

// ---------- getIPAvailability ----------

type NeoliteIPAvailability struct{}

type NeoliteIPAvailabilityArgs struct {
	ProductID int64 `pulumi:"productId"`
}

type NeoliteIPAvailabilityResult struct {
	Available bool `pulumi:"available"`
}

func (a *NeoliteIPAvailabilityArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id NEO Lite.")
}

func (NeoliteIPAvailability) Invoke(
	ctx context.Context, req infer.FunctionRequest[NeoliteIPAvailabilityArgs],
) (infer.FunctionResponse[NeoliteIPAvailabilityResult], error) {
	c := GetClient(ctx)
	out, err := c.Neolite().ProductIPAvailability(ctx, req.Input.ProductID)
	if err != nil {
		return infer.FunctionResponse[NeoliteIPAvailabilityResult]{}, err
	}
	return infer.FunctionResponse[NeoliteIPAvailabilityResult]{
		Output: NeoliteIPAvailabilityResult{Available: out.Available},
	}, nil
}
