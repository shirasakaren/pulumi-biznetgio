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
