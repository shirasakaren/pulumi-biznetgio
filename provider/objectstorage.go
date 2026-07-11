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
