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
					fl, ok := raw.(map[string]any)
					if !ok {
						continue
					}
					flavor := GpuFlavor{Raw: string(gpuJSON(fl))}
					if v, ok := gpuInt64(fl, "flavor_id", "id", "product_id"); ok {
						flavor.FlavorID = v
					}
					flavor.Name = gpuStringDefault(fl, "name", "flavor_name", "label")
					product.Flavors = append(product.Flavors, flavor)
				}
			}
		}
		out = append(out, product)
	}
	return infer.FunctionResponse[GpuProductsState]{Output: GpuProductsState{Products: out}}, nil
}

type GpuConsole struct{}

type GpuConsoleArgs struct {
	AccountID string `pulumi:"accountId"`
}

type GpuConsoleResult struct {
	URL string `pulumi:"url" provider:"secret"`
	Raw string `pulumi:"raw" provider:"secret"`
}

func (a *GpuConsoleArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AccountID, "Account id of the GPU instance to get console access for.")
}

func (r *GpuConsoleResult) Annotate(ann infer.Annotator) {
	ann.Describe(&r.URL, "Console access URL from `POST /neo-gpus/accounts/{account_id}/console-access`. "+
		"Side-effecting: each call may mint a new one-time session, do not use inside plan diffing.")
	ann.Describe(&r.Raw, "Raw JSON of the full response.")
}

func (GpuConsole) Invoke(
	ctx context.Context, req infer.FunctionRequest[GpuConsoleArgs],
) (infer.FunctionResponse[GpuConsoleResult], error) {
	c := GetClient(ctx)
	accountID, err := strconv.ParseInt(req.Input.AccountID, 10, 64)
	if err != nil {
		return infer.FunctionResponse[GpuConsoleResult]{},
			fmt.Errorf("biznetgio: invalid accountId %q: %s", req.Input.AccountID, err)
	}
	m, err := c.GPU().ConsoleAccess(ctx, accountID)
	if err != nil {
		return infer.FunctionResponse[GpuConsoleResult]{}, err
	}
	u, ok := gpuString(m, "url", "console_url", "access_url", "href")
	if !ok {
		return infer.FunctionResponse[GpuConsoleResult]{},
			fmt.Errorf("biznetgio: console access response missing url: %s", gpuJSON(m))
	}
	return infer.FunctionResponse[GpuConsoleResult]{Output: GpuConsoleResult{
		URL: u,
		Raw: string(gpuJSON(m)),
	}}, nil
}

type GpuGraph struct{}

type GpuGraphArgs struct {
	AccountID string  `pulumi:"accountId"`
	Timeframe *string `pulumi:"timeframe,optional"`
}

type GpuGraphResult struct {
	Graph string `pulumi:"graph"`
	Raw   string `pulumi:"raw" provider:"secret"`
}

func (a *GpuGraphArgs) Annotate(ann infer.Annotator) {
