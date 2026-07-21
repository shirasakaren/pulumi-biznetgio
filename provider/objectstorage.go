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

type ObjectStorage struct{}

type ObjectStorageArgs struct {
	ProductID         int64   `pulumi:"productId"`
	Cycle             string  `pulumi:"cycle"`
	Label             string  `pulumi:"label" provider:"replaceOnChanges"`
	Quota             *int64  `pulumi:"quota,optional"`
	Promocode         *string `pulumi:"promocode,optional" provider:"replaceOnChanges"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
}

type ObjectStorageState struct {
	ObjectStorageArgs
	OrderID *string `pulumi:"orderId"`
	Status  string  `pulumi:"status"`
	Raw     string  `pulumi:"raw" provider:"secret"`
}

func (a *ObjectStorageArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Object storage product/plan id.")
	ann.Describe(&a.Cycle, "Billing cycle, e.g. `m` for monthly or `a` for annual.")
	ann.Describe(&a.Label, "Instance label, 6-16 chars, `[a-zA-Z0-9-_]`. Create-only, changing it replaces the instance.")
	ann.Describe(&a.Quota, "Quota in GB. Defaults to 10. "+
		"Hanya bisa diperbesar (dihitung sebagai tambahan dari quota lama).")
	ann.SetDefault(&a.Quota, int64(10))
	ann.Describe(&a.Promocode, "Promo code to apply at creation. Create-only.")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card. Defaults to true.")
	ann.SetDefault(&a.PayWithCreditCard, true)
}

func (s *ObjectStorageState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id from the creation response.")
	ann.Describe(&s.Status, "Current status of the instance: Active, Pending, Suspended, or Terminated.")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (ObjectStorage) Create(
	ctx context.Context, req infer.CreateRequest[ObjectStorageArgs],
) (infer.CreateResponse[ObjectStorageState], error) {
	resp := infer.CreateResponse[ObjectStorageState]{Output: osStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs

	raw, err := c.ObjectStorage().Create(ctx, client.NOSCreatePayload{
		ProductID:        a.ProductID,
		Cycle:            a.Cycle,
		Label:            a.Label,
		Quota:            osQuota(a),
		Promocode:        osStr(a.Promocode),
		PayInvoiceWithCC: osPay(a),
	})
	if err != nil {
		return infer.CreateResponse[ObjectStorageState]{}, err
	}
	accountID, ok := osString(raw, "account_id", "id")
	if !ok {
		return infer.CreateResponse[ObjectStorageState]{},
			fmt.Errorf("biznetgio: object storage create response missing account_id: %s", osJSON(raw))
	}
	resp.ID = accountID
	resp.Output = osStateFromMap(ctx, a, raw)

	if _, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return c.ObjectStorage().AccountGet(ctx, osParseID(accountID))
		},
		osStatus, []string{"active"}, []string{"terminated", "suspended", "failed", "error", "deleted",
			"cancelled"}); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return resp, infer.ResourceInitFailedError{Reasons: []string{
				fmt.Sprintf("object storage %s belum active, lanjutin via update aja: %s", accountID, err),
			}}
		}
		return resp, err
	}
	return resp, nil
}

func (ObjectStorage) Update(
	ctx context.Context, req infer.UpdateRequest[ObjectStorageArgs, ObjectStorageState],
) (infer.UpdateResponse[ObjectStorageState], error) {
	resp := infer.UpdateResponse[ObjectStorageState]{Output: osStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs
	accountID := osParseID(req.ID)

	old := osQuota(req.State.ObjectStorageArgs)
	cur := osQuota(a)
	if cur < old {
		return infer.UpdateResponse[ObjectStorageState]{},
			fmt.Errorf("biznetgio: object storage quota hanya bisa diperbesar, bukan diperkecil (%d -> %d)", old, cur)
	}
	if cur > old {
		if _, err := c.ObjectStorage().QuotaUpgrade(ctx, accountID, client.NOSQuotaUpgradePayload{
			AddQuota:         cur - old,
			PayInvoiceWithCC: osPay(a),
		}); err != nil {
			return infer.UpdateResponse[ObjectStorageState]{}, err
		}

		final, err := client.WaitForStatus(ctx, 5*time.Second,
			func(ctx context.Context) (map[string]any, error) { return c.ObjectStorage().AccountGet(ctx, accountID) },
			osStatus, []string{"active"}, []string{"terminated", "suspended", "failed", "error", "deleted", "cancelled"})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return infer.UpdateResponse[ObjectStorageState]{
					Output: osStateFromMap(ctx, a, nil),
				}, infer.ResourceInitFailedError{Reasons: []string{err.Error()}}
			}
			return infer.UpdateResponse[ObjectStorageState]{}, err
		}
		resp.Output = osStateFromMap(ctx, a, final)
	}
	return resp, nil
}

func (ObjectStorage) Read(
	ctx context.Context, req infer.ReadRequest[ObjectStorageArgs, ObjectStorageState],
) (infer.ReadResponse[ObjectStorageArgs, ObjectStorageState], error) {
	resp := infer.ReadResponse[ObjectStorageArgs, ObjectStorageState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  osStateFromMap(ctx, req.Inputs, nil),
	}
	c := GetClient(ctx)
	m, err := c.ObjectStorage().AccountGet(ctx, osParseID(req.ID))
	if err != nil {
		if client.IsNotFound(err) {
			return resp, fmt.Errorf("biznetgio: object storage %s not found", req.ID)
		}
		return resp, err
	}
	if strings.EqualFold(osStringDefault(m, "status", "state"), "terminated") {
		return resp, fmt.Errorf("biznetgio: object storage %s terminated", req.ID)
	}
	resp.State = osStateFromMap(ctx, req.Inputs, m)
	return resp, nil
}

func (ObjectStorage) Delete(
	ctx context.Context, req infer.DeleteRequest[ObjectStorageState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	if _, err := c.ObjectStorage().Delete(ctx, osParseID(req.ID)); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func osPay(a ObjectStorageArgs) string {
	if a.PayWithCreditCard != nil && !*a.PayWithCreditCard {
		return "no"
	}
	return "yes"
}

func osQuota(a ObjectStorageArgs) int64 {
	if a.Quota != nil {
		return *a.Quota
	}
	return 10
}

func osParseID(id string) int64 {
	n, _ := strconv.ParseInt(id, 10, 64)
	return n
}

func osStateFromMap(
	_ context.Context, args ObjectStorageArgs, m map[string]any,
) ObjectStorageState {
	st := ObjectStorageState{ObjectStorageArgs: args}
	if m == nil {
		return st
	}
	if v, ok := osInt(m, "product_id"); ok {
		st.ProductID = v
	}
	if v, ok := osString(m, "cycle", "billingcycle"); ok {
		st.Cycle = v
	}
	if v, ok := osString(m, "label", "name", "service_name"); ok {
		st.Label = v
	}
	if v, ok := osInt(m, "quota"); ok {
		st.Quota = &v
	}
	if v, ok := osString(m, "order_id"); ok {
		st.OrderID = &v
	}
	st.Status = osStringDefault(m, "status", "state")
	st.Raw = string(osJSON(m))
	return st
