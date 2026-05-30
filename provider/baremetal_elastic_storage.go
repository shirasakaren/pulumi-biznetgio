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

type BaremetalElasticStorage struct{}

type BaremetalElasticStorageArgs struct {
	ProductID         int64   `pulumi:"productId"`
	Cycle             string  `pulumi:"cycle"`
	StorageName       string  `pulumi:"storageName"`
	MetalAccountID    int64   `pulumi:"metalAccountId"`
	Size              *int64  `pulumi:"size,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
}

type BaremetalElasticStorageState struct {
	BaremetalElasticStorageArgs
	Status    string  `pulumi:"status"`
	CreatedAt *string `pulumi:"createdAt"`
	Raw       string  `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalElasticStorageArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id from `GET /baremetal-neo-elastic-storages/products`. "+
		"Changing it triggers change-package (`POST .../{account_id}`).")
	ann.Describe(&a.Cycle, "Billing cycle, e.g. `m` for monthly or `a` for annual. "+
		"Create-only, changing it replaces the storage.")
	ann.Describe(&a.StorageName, "Name of the storage. Create-only, changing it replaces the storage.")
	ann.Describe(&a.MetalAccountID, "Account id of the target baremetal. "+
		"Can only be set at creation (no re-attach endpoint).")
	ann.Describe(&a.Size, "Storage size in GB. Defaults to 100. Changing it triggers upgrade "+
		"(`PUT .../{account_id}`) — grow-only, shrinking is rejected by the API.")
	ann.SetDefault(&a.Size, int64(100))
	ann.Describe(&a.Promocode, "Promo code to apply at creation.")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card. Defaults to true. "+
		"Set false to leave the invoice unpaid in the portal; the resource stays Pending until paid.")
	ann.SetDefault(&a.PayWithCreditCard, true)
}

func (s *BaremetalElasticStorageState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Status, "Current status of the storage (e.g. Active, Pending, Suspended, Terminated).")
	ann.Describe(&s.CreatedAt, "Creation timestamp (alias of created_at/date_created).")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (BaremetalElasticStorage) Create(
	ctx context.Context, req infer.CreateRequest[BaremetalElasticStorageArgs],
) (infer.CreateResponse[BaremetalElasticStorageState], error) {
	resp := infer.CreateResponse[BaremetalElasticStorageState]{Output: elasticStorageStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs

	pay := "yes"
	if a.PayWithCreditCard != nil && !*a.PayWithCreditCard {
		pay = "no"
	}
	raw, err := c.BaremetalElasticStorage().Create(ctx, client.NeoElasticStorageCreatePayload{
		ProductID:        a.ProductID,
		Cycle:            a.Cycle,
		StorageName:      a.StorageName,
		MetalAccountID:   a.MetalAccountID,
		Size:             bmtInt64Ptr(a.Size),
		Promocode:        bmtStr(a.Promocode),
		PayInvoiceWithCC: pay,
	})
	if err != nil {
		return infer.CreateResponse[BaremetalElasticStorageState]{}, err
	}
	accountID, ok := bmtInt64(raw, "account_id", "id")
	if !ok {
		return infer.CreateResponse[BaremetalElasticStorageState]{},
			fmt.Errorf("biznetgio: elastic storage create response missing account_id: %s", bmtJSON(raw))
	}
	resp.ID = strconv.FormatInt(accountID, 10)
	resp.Output = elasticStorageStateFromMap(ctx, a, raw)

	final, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return c.BaremetalElasticStorage().Get(ctx, accountID)
		},
		bmtStatus, []string{"active"}, []string{"terminated", "error", "failed"})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return resp, infer.ResourceInitFailedError{Reasons: []string{
				fmt.Sprintf("elastic storage %d belum active, lanjutin via update aja: %s", accountID, err),
			}}
		}
		return resp, err
	}
	resp.Output = elasticStorageStateFromMap(ctx, a, final)
	return resp, nil
}

func (BaremetalElasticStorage) Update(
	ctx context.Context, req infer.UpdateRequest[BaremetalElasticStorageArgs, BaremetalElasticStorageState],
) (infer.UpdateResponse[BaremetalElasticStorageState], error) {
	resp := infer.UpdateResponse[BaremetalElasticStorageState]{Output: elasticStorageStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.UpdateResponse[BaremetalElasticStorageState]{},
			fmt.Errorf("biznetgio: invalid elastic storage id %q: %s", req.ID, err)
	}

	pay := "yes"
	if a.PayWithCreditCard != nil && !*a.PayWithCreditCard {
		pay = "no"
	}
	changed := false
	if bmtInt64Ptr(a.Size) != bmtInt64Ptr(req.State.Size) {
		if _, err := c.BaremetalElasticStorage().Upgrade(ctx, accountID, client.UpgradeNeoElasticStorage{
			Size:             bmtInt64Ptr(a.Size),
			PayInvoiceWithCC: pay,
		}); err != nil {
			return infer.UpdateResponse[BaremetalElasticStorageState]{},
				fmt.Errorf("biznetgio: upgrade elastic storage %d size: %w", accountID, err)
		}
		changed = true
	}
	if a.ProductID != req.State.ProductID {
		if _, err := c.BaremetalElasticStorage().ChangePackage(ctx, accountID, client.ChangePackageNeoElasticStorage{
			NewProductID:     a.ProductID,
			PayInvoiceWithCC: pay,
		}); err != nil {
			return infer.UpdateResponse[BaremetalElasticStorageState]{},
				fmt.Errorf("biznetgio: change elastic storage %d package: %w", accountID, err)
		}
		changed = true
	}

	if changed {
		final, err := client.WaitForStatus(ctx, 5*time.Second,
			func(ctx context.Context) (map[string]any, error) {
				return c.BaremetalElasticStorage().Get(ctx, accountID)
			},
			bmtStatus, []string{"active"}, []string{"terminated", "error", "failed"})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return infer.UpdateResponse[BaremetalElasticStorageState]{
					Output: elasticStorageStateFromMap(ctx, a, nil),
				}, infer.ResourceInitFailedError{Reasons: []string{err.Error()}}
			}
			return infer.UpdateResponse[BaremetalElasticStorageState]{}, err
		}
		resp.Output = elasticStorageStateFromMap(ctx, a, final)
	}
	return resp, nil
}

func (BaremetalElasticStorage) Read(
	ctx context.Context, req infer.ReadRequest[BaremetalElasticStorageArgs, BaremetalElasticStorageState],
) (infer.ReadResponse[BaremetalElasticStorageArgs, BaremetalElasticStorageState], error) {
	resp := infer.ReadResponse[BaremetalElasticStorageArgs, BaremetalElasticStorageState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  elasticStorageStateFromMap(ctx, req.Inputs, nil),
	}
	c := GetClient(ctx)
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return resp, fmt.Errorf("biznetgio: invalid elastic storage id %q: %s", req.ID, err)
	}
