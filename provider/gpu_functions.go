package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type GpuProducts struct{}

type GpuProductsArgs struct{}

type GpuProductsState struct {
	Products []GpuProduct `pulumi:"products"`
}

type GpuProduct struct {
	ProductID    int64       `pulumi:"productId"`
	Name         string      `pulumi:"name"`
	Description  string      `pulumi:"description"`
	CategoryName string      `pulumi:"categoryName"`
	Flavors      []GpuFlavor `pulumi:"flavors"`
	Raw          string      `pulumi:"raw" provider:"secret"`
}

type GpuFlavor struct {
	FlavorID int64  `pulumi:"flavorId"`
	Name     string `pulumi:"name"`
	Raw      string `pulumi:"raw" provider:"secret"`
}

func (a *GpuProductsArgs) Annotate(_ infer.Annotator) {
}

func (s *GpuProductsState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Products, "List of GPU products available for order.")
