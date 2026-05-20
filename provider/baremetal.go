package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type Baremetal struct{}

type BaremetalArgs struct {
	ProductID         int64   `pulumi:"productId" provider:"replaceOnChanges"`
	Cycle             string  `pulumi:"cycle" provider:"replaceOnChanges"`
	SelectOS          *string `pulumi:"selectOs,optional" provider:"replaceOnChanges"`
	KeypairID         int64   `pulumi:"keypairId" provider:"replaceOnChanges"`
	Label             string  `pulumi:"label"`
	PublicIP          *int64  `pulumi:"publicIp,optional" provider:"replaceOnChanges"`
	Promocode         *string `pulumi:"promocode,optional" provider:"replaceOnChanges"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional" provider:"replaceOnChanges"`
	PowerState        *string `pulumi:"powerState,optional"`
	ResetTrigger      *string `pulumi:"resetTrigger,optional"`
	RebuildOS         *string `pulumi:"rebuildOs,optional"`
}

type BaremetalState struct {
	BaremetalArgs
	Status    string  `pulumi:"status"`
	OrderID   *string `pulumi:"orderId"`
	IPAddress *string `pulumi:"ipAddress"`
	CreatedAt *string `pulumi:"createdAt"`
	Raw       string  `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id from the baremetal products function.")
	ann.Describe(&a.Cycle, "Billing cycle, e.g. `m` for monthly or `a` for annual.")
	ann.Describe(&a.SelectOS, "OS to install at creation, from `GET /baremetals/products/{product_id}/oss`. "+
		"Defaults to `ubuntu-22`. Create-only, changing it replaces the instance.")
	ann.SetDefault(&a.SelectOS, "ubuntu-22")
	ann.Describe(&a.KeypairID, "Baremetal keypair id, from BaremetalKeypair. "+
		"The baremetal keypair pool is separate from neolite/gpu.")
	ann.Describe(&a.Label, "Display name of the server. The only field updatable in place.")
	ann.Describe(&a.PublicIP, "Number of public IPs to request (1 = with public IP). Defaults to 1. "+
		"Create-only, changing it replaces the instance.")
	ann.SetDefault(&a.PublicIP, int64(1))
	ann.Describe(&a.Promocode, "Promo code to apply at creation.")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card. Defaults to true. "+
		"Set false to leave the invoice unpaid in the portal; the resource stays Pending until paid.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.PowerState, "Power state of the server: `on` or `off`. Only fires an API call when the value changes.")
	ann.Describe(&a.ResetTrigger, "Change this value to trigger a one-shot reset/reboot. "+
		"The reset is an action, not a stable state.")
	ann.Describe(&a.RebuildOS, "Change this value to rebuild the instance (destructive, wipes the OS) "+
		"via `PUT /baremetals/{account_id}/rebuild`. Valid OS list comes from the rebuildOsList function.")
}

func (s *BaremetalState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Status, "Current status of the server.")
	ann.Describe(&s.OrderID, "Order id from the creation response.")
	ann.Describe(&s.IPAddress, "Public IP address of the server, when present.")
	ann.Describe(&s.CreatedAt, "Creation timestamp (alias of created_at/date_created).")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (Baremetal) Create(
	ctx context.Context, req infer.CreateRequest[BaremetalArgs],
) (infer.CreateResponse[BaremetalState], error) {
	resp := infer.CreateResponse[BaremetalState]{Output: baremetalStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs

	pay := "yes"
	if a.PayWithCreditCard != nil && !*a.PayWithCreditCard {
		pay = "no"
	}
	raw, err := c.Baremetal().Create(ctx, client.BaremetalCreatePayload{
		ProductID:        a.ProductID,
		Cycle:            a.Cycle,
		SelectOS:         bmtStr(a.SelectOS),
		KeypairID:        a.KeypairID,
		Label:            a.Label,
		PublicIP:         bmtInt64Ptr(a.PublicIP),
		Promocode:        bmtStr(a.Promocode),
		PayInvoiceWithCC: pay,
	})
	if err != nil {
		return infer.CreateResponse[BaremetalState]{}, err
	}
	accountID, ok := bmtInt64(raw, "account_id", "id")
	if !ok {
		return infer.CreateResponse[BaremetalState]{},
			fmt.Errorf("biznetgio: baremetal create response missing account_id: %s", bmtJSON(raw))
	}
	resp.ID = strconv.FormatInt(accountID, 10)
	resp.Output = baremetalStateFromMap(ctx, a, raw)

	final, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) { return c.Baremetal().AccountGet(ctx, accountID) },
		bmtStatus, []string{"active"}, []string{"terminated", "error", "failed"})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return resp, infer.ResourceInitFailedError{Reasons: []string{
				fmt.Sprintf("baremetal %d belum active, lanjutin via update aja: %s", accountID, err),
			}}
		}
		return resp, err
	}
	resp.Output = baremetalStateFromMap(ctx, a, final)
	if st, err := c.Baremetal().StateGet(ctx, accountID); err == nil {
		if v, ok := bmtString(st, "state", "status"); ok {
			resp.Output.PowerState = &v
		}
	}
	return resp, nil
}

func (Baremetal) Update(ctx context.Context,
	req infer.UpdateRequest[BaremetalArgs, BaremetalState],
) (infer.UpdateResponse[BaremetalState], error) {
	resp := infer.UpdateResponse[BaremetalState]{Output: baremetalStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.UpdateResponse[BaremetalState]{}, fmt.Errorf("biznetgio: invalid baremetal id %q: %s", req.ID, err)
	}

	if a.Label != req.State.Label {
		if _, err := c.Baremetal().UpdateLabel(ctx, accountID, a.Label); err != nil {
			return infer.UpdateResponse[BaremetalState]{}, fmt.Errorf("biznetgio: update baremetal label %d: %w", accountID, err)
		}
	}
	if a.PowerState != nil && *a.PowerState != bmtStr(req.State.PowerState) {
		if _, err := c.Baremetal().StateSet(ctx, accountID, *a.PowerState); err != nil {
			return infer.UpdateResponse[BaremetalState]{},
				fmt.Errorf("biznetgio: set baremetal power state %q on %d: %w", *a.PowerState, accountID, err)
		}
	}
	if t := bmtStr(a.ResetTrigger); t != "" && t != bmtStr(req.State.ResetTrigger) {
		if _, err := c.Baremetal().StateSet(ctx, accountID, "reset"); err != nil {
			return infer.UpdateResponse[BaremetalState]{}, fmt.Errorf("biznetgio: reset baremetal %d: %w", accountID, err)
		}
	}
	if t := bmtStr(a.RebuildOS); t != "" && t != bmtStr(req.State.RebuildOS) {
		if _, err := c.Baremetal().Rebuild(ctx, accountID, t); err != nil {
			return infer.UpdateResponse[BaremetalState]{}, fmt.Errorf("biznetgio: rebuild baremetal %d: %w", accountID, err)
		}
		final, err := client.WaitForStatus(ctx, 5*time.Second,
			func(ctx context.Context) (map[string]any, error) { return c.Baremetal().AccountGet(ctx, accountID) },
			bmtStatus, []string{"active"}, []string{"terminated", "error", "failed"})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return infer.UpdateResponse[BaremetalState]{Output: baremetalStateFromMap(ctx, a, nil)},
					infer.ResourceInitFailedError{Reasons: []string{err.Error()}}
			}
			return infer.UpdateResponse[BaremetalState]{}, err
		}
		resp.Output = baremetalStateFromMap(ctx, a, final)
	}
	return resp, nil
}

func (Baremetal) Read(ctx context.Context,
	req infer.ReadRequest[BaremetalArgs, BaremetalState],
) (infer.ReadResponse[BaremetalArgs, BaremetalState], error) {
	resp := infer.ReadResponse[BaremetalArgs, BaremetalState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  baremetalStateFromMap(ctx, req.Inputs, nil),
	}
	c := GetClient(ctx)
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return resp, fmt.Errorf("biznetgio: invalid baremetal id %q: %s", req.ID, err)
	}
	m, err := c.Baremetal().AccountGet(ctx, accountID)
	if err != nil {
		if client.IsNotFound(err) {
			return resp, fmt.Errorf("biznetgio: baremetal %s not found", req.ID)
		}
		return resp, err
	}
	resp.State = baremetalStateFromMap(ctx, req.Inputs, m)
	if st, err := c.Baremetal().StateGet(ctx, accountID); err == nil {
		if v, ok := bmtString(st, "state", "status"); ok {
			resp.State.PowerState = &v
		}
	}
	return resp, nil
}

func (Baremetal) Delete(ctx context.Context, req infer.DeleteRequest[BaremetalState]) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("biznetgio: invalid baremetal id %q: %s", req.ID, err)
	}
	if _, err := c.Baremetal().Delete(ctx, accountID); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func baremetalStateFromMap(_ context.Context, args BaremetalArgs, m map[string]any) BaremetalState {
	st := BaremetalState{BaremetalArgs: args}
	if m == nil {
		return st
	}
	if v, ok := bmtInt64(m, "product_id"); ok {
		st.ProductID = v
	}
	if v, ok := bmtInt64(m, "keypair_id", "neosshkey_id"); ok {
		st.KeypairID = v
	}
	if v, ok := bmtString(m, "select_os", "os"); ok {
		st.SelectOS = &v
	}
	if v, ok := bmtString(m, "label", "name"); ok {
		st.Label = v
	}
	if v, ok := bmtString(m, "order_id"); ok {
		st.OrderID = &v
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

func bmtStatus(m map[string]any) string {
	return strings.ToLower(bmtStringDefault(m, "status", "state"))
}

func bmtStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func bmtInt64Ptr(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func bmtInt64(m map[string]any, keys ...string) (int64, bool) {
	for _, k := range keys {
		v, ok := m[k]
		if !ok {
			continue
		}
		switch n := v.(type) {
		case float64:
			return int64(n), true
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return i, true
			}
		case string:
