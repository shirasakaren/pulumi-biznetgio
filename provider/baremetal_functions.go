package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// BaremetalProduct mirrors one item from `GET /baremetals/products`.
type BaremetalProduct struct {
	ProductID   int64  `pulumi:"productId"`
	Name        string `pulumi:"name"`
	Description string `pulumi:"description"`
	Raw         string `pulumi:"raw" provider:"secret"`
}

func (p *BaremetalProduct) Annotate(ann infer.Annotator) {
	ann.Describe(&p.ProductID, "Product id, for the `productId` input of Baremetal.")
	ann.Describe(&p.Name, "Product name (alias of name/product_name/label).")
	ann.Describe(&p.Description, "Product description, when present.")
	ann.Describe(&p.Raw, "Raw JSON of the product item, for anything not modeled yet.")
}

type BaremetalProducts struct{}

type BaremetalProductsArgs struct{}

type BaremetalProductsResult struct {
	Products []BaremetalProduct `pulumi:"products"`
	Raw      string             `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalProductsArgs) Annotate(_ infer.Annotator) {}

func (r *BaremetalProductsResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.Products, "Available baremetal products from `GET /baremetals/products`.")
	ann.Describe(&r.Raw, "Raw JSON of the full list response.")
}

func (BaremetalProducts) Invoke(
	ctx context.Context, _ infer.FunctionRequest[BaremetalProductsArgs],
) (infer.FunctionResponse[BaremetalProductsResult], error) {
	c := GetClient(ctx)
	items, err := c.Baremetal().ProductList(ctx)
	if err != nil {
		return infer.FunctionResponse[BaremetalProductsResult]{}, err
	}
	out := BaremetalProductsResult{Products: make([]BaremetalProduct, 0, len(items))}
	for _, it := range items {
		out.Products = append(out.Products, BaremetalProduct{
			ProductID:   bmtStringDefaultInt64(it, "product_id", "id"),
			Name:        bmtStringDefault(it, "name", "product_name", "label"),
			Description: bmtStringDefault(it, "description"),
			Raw:         string(bmtJSON(it)),
		})
	}
	out.Raw = string(bmtJSONList(items))
	return infer.FunctionResponse[BaremetalProductsResult]{Output: out}, nil
}

type BaremetalRebuildOsList struct{}

type BaremetalRebuildOsListArgs struct {
	AccountID int64 `pulumi:"accountId"`
}

type BaremetalRebuildOsListResult struct {
	Oss []string `pulumi:"oss"`
	Raw string   `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalRebuildOsListArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Account id of the baremetal to list rebuild OS options for.")
}

func (r *BaremetalRebuildOsListResult) Annotate(ann infer.Annotator) {
