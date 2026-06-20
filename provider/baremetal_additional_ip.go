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
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card. Defaults to true. "+
		"Set false to leave the invoice unpaid in the portal; the resource stays Pending until paid.")
	ann.SetDefault(&a.PayWithCreditCard, true)
}

func (s *BaremetalAdditionalIpState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Status, "Current status of the IP (e.g. Active, Pending, Suspended, Terminated).")
	ann.Describe(&s.IPAddress, "Assigned IP address, when present in the response.")
	ann.Describe(&s.CreatedAt, "Creation timestamp (alias of created_at/date_created).")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (BaremetalAdditionalIp) Create(
	ctx context.Context, req infer.CreateRequest[BaremetalAdditionalIpArgs],
) (infer.CreateResponse[BaremetalAdditionalIpState], error) {
	resp := infer.CreateResponse[BaremetalAdditionalIpState]{Output: additionalIpStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs

	pay := "yes"
	if a.PayWithCreditCard != nil && !*a.PayWithCreditCard {
		pay = "no"
	}
	raw, err := c.BaremetalAdditionalIP().Create(ctx, client.AdditionalIPCreatePayload{
		ProductID:        a.ProductID,
		Cycle:            a.Cycle,
		Region:           bmtStr(a.Region),
		Promocode:        bmtStr(a.Promocode),
		PayInvoiceWithCC: pay,
	})
	if err != nil {
		return infer.CreateResponse[BaremetalAdditionalIpState]{}, err
	}
	accountID, ok := bmtInt64(raw, "account_id", "id")
	if !ok {
		return infer.CreateResponse[BaremetalAdditionalIpState]{},
