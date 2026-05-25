package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type BaremetalAdditionalIp struct{}

type BaremetalAdditionalIpArgs struct {
	ProductID         int64   `pulumi:"productId"`
	Cycle             string  `pulumi:"cycle"`
	Region            *string `pulumi:"region,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
}

type BaremetalAdditionalIpState struct {
	BaremetalAdditionalIpArgs
	Status    string  `pulumi:"status"`
	IPAddress *string `pulumi:"ipAddress"`
	CreatedAt *string `pulumi:"createdAt"`
	Raw       string  `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalAdditionalIpArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id from `GET /baremetal-additional-ips/products`.")
	ann.Describe(&a.Cycle, "Billing cycle, e.g. `m` for monthly or `a` for annual.")
	ann.Describe(&a.Region, "Datacenter region; valid list from `GET /baremetal-additional-ips/regions`. "+
		"Defaults to `wjv-1`. Create-only, changing it replaces the IP.")
	ann.SetDefault(&a.Region, "wjv-1")
	ann.Describe(&a.Promocode, "Promo code to apply at creation.")
// wip 178
