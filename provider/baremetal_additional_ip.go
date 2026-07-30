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
			fmt.Errorf("biznetgio: additional ip create response missing account_id: %s", bmtJSON(raw))
	}
	resp.ID = strconv.FormatInt(accountID, 10)
	resp.Output = additionalIpStateFromMap(ctx, a, raw)

	final, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return c.BaremetalAdditionalIP().Get(ctx, strconv.FormatInt(accountID, 10))
		},
		bmtStatus, []string{"active"}, []string{"terminated", "error", "failed"})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return resp, infer.ResourceInitFailedError{Reasons: []string{
				fmt.Sprintf("additional ip %d belum active, lanjutin via update aja: %s", accountID, err),
			}}
		}
		return resp, err
	}
	resp.Output = additionalIpStateFromMap(ctx, a, final)
	return resp, nil
}

func (BaremetalAdditionalIp) Read(
	ctx context.Context, req infer.ReadRequest[BaremetalAdditionalIpArgs, BaremetalAdditionalIpState],
) (infer.ReadResponse[BaremetalAdditionalIpArgs, BaremetalAdditionalIpState], error) {
	resp := infer.ReadResponse[BaremetalAdditionalIpArgs, BaremetalAdditionalIpState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  additionalIpStateFromMap(ctx, req.Inputs, nil),
	}
	c := GetClient(ctx)
	m, err := c.BaremetalAdditionalIP().Get(ctx, req.ID)
	if err != nil {
		if client.IsNotFound(err) {
			return resp, fmt.Errorf("biznetgio: additional ip %s not found", req.ID)
		}
		return resp, err
	}
	resp.State = additionalIpStateFromMap(ctx, req.Inputs, m)
	return resp, nil
}

func (BaremetalAdditionalIp) Delete(
	ctx context.Context, req infer.DeleteRequest[BaremetalAdditionalIpState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("biznetgio: invalid additional ip id %q: %s", req.ID, err)
	}
	if _, err := c.BaremetalAdditionalIP().Delete(ctx, accountID); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func additionalIpStateFromMap(
	_ context.Context, args BaremetalAdditionalIpArgs, m map[string]any,
) BaremetalAdditionalIpState {
	st := BaremetalAdditionalIpState{BaremetalAdditionalIpArgs: args}
	if m == nil {
		return st
	}
	if v, ok := bmtInt64(m, "product_id"); ok {
		st.ProductID = v
	}
	if v, ok := bmtString(m, "ip", "public_ip", "ip_address", "ipv4"); ok {
		st.IPAddress = &v
	}
	if v, ok := bmtString(m, "created_at", "date_created"); ok {
		st.CreatedAt = &v
	}
	st.Status = bmtStringDefault(m, "status", "state")
	st.Raw = string(bmtJSON(m))
	return st
}
