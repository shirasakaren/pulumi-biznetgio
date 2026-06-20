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
