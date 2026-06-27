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
}

func (p *GpuProduct) Annotate(ann infer.Annotator) {
	ann.Describe(&p.ProductID, "Product id, use it as GpuInstance productId.")
	ann.Describe(&p.Name, "Product name.")
	ann.Describe(&p.Description, "Product description.")
	ann.Describe(&p.CategoryName, "Category of the product.")
	ann.Describe(&p.Flavors, "Available flavors for this product.")
	ann.Describe(&p.Raw, "Raw JSON of the product response.")
}

func (f *GpuFlavor) Annotate(ann infer.Annotator) {
	ann.Describe(&f.FlavorID, "Flavor id.")
	ann.Describe(&f.Name, "Flavor name.")
	ann.Describe(&f.Raw, "Raw JSON of the flavor response.")
}

func (GpuProducts) Invoke(
	ctx context.Context, _ infer.FunctionRequest[GpuProductsArgs],
) (infer.FunctionResponse[GpuProductsState], error) {
	c := GetClient(ctx)
	products, err := c.GPU().ProductList(ctx)
	if err != nil {
		return infer.FunctionResponse[GpuProductsState]{}, err
	}
	out := make([]GpuProduct, 0, len(products))
	for _, m := range products {
		product := GpuProduct{Raw: string(gpuJSON(m))}
		if v, ok := gpuInt64(m, "product_id", "id"); ok {
			product.ProductID = v
		}
		product.Name = gpuStringDefault(m, "name", "product_name", "label")
		product.Description = gpuStringDefault(m, "description")
		product.CategoryName = gpuStringDefault(m, "category_name", "category")
		if product.ProductID != 0 {
			flavors, err := c.GPU().ProductFlavors(ctx, product.ProductID)
			if err == nil {
				for _, raw := range flavors {
