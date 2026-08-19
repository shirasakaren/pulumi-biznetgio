package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
)

type NeoliteProDisk struct{}

type NeoliteProDiskArgs struct {
	ProductID         int64   `pulumi:"productId"`
	Cycle             string  `pulumi:"cycle"`
	NeoliteAccountID  int64   `pulumi:"neoliteAccountId"`
	ServiceName       *string `pulumi:"serviceName,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
	Size              *int64  `pulumi:"size,optional"`
}

type NeoliteProDiskState struct {
	NeoliteProDiskArgs
	OrderID string `pulumi:"orderId"`
	Status  string `pulumi:"status"`
	Raw     string `pulumi:"raw" provider:"secret"`
}

func (a *NeoliteProDiskArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id disk dari endpoint `/neolite-pros/disks/products`.")
	ann.Describe(&a.Cycle, "Siklus billing: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.NeoliteAccountID, "Account id VM NEO Lite Pro tempat disk dipasang.")
	ann.Describe(&a.ServiceName, "Nama layanan disk. Default `service-name`.")
	ann.SetDefault(&a.ServiceName, "service-name")
	ann.Describe(&a.Promocode, "Kode promo saat order.")
	ann.SetDefault(&a.Promocode, "")
	ann.Describe(&a.PayWithCreditCard, "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). "+
		"Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.Size, "Ukuran disk (GB). Default 30 (minimal product disk pro). Cuma bisa naik, bukan turun.")
	ann.SetDefault(&a.Size, int64(30))
}

func (s *NeoliteProDiskState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id dari response create.")
	ann.Describe(&s.Status, "Status disk (Active, Pending, Suspended, Terminated).")
	ann.Describe(&s.Raw, "Full JSON response disk terakhir dari API, buat akses field yang belum dimodel.")
}

func (NeoliteProDisk) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteProDiskArgs],
) (infer.CreateResponse[NeoliteProDiskState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteProDiskState]{
			ID:     "0",
			Output: NeoliteProDiskState{NeoliteProDiskArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	out, err := c.NeolitePro().DiskCreate(ctx, client.NeoliteDiskCreatePayload{
		ProductID:        in.ProductID,
		Cycle:            in.Cycle,
		NeoliteAccountID: in.NeoliteAccountID,
		ServiceName:      strPtr(in.ServiceName),
		Promocode:        strPtr(in.Promocode),
		PayInvoiceWithCC: yesNo(in.PayWithCreditCard),
		Size:             i64Val(in.Size),
	})
	if err != nil {
		return infer.CreateResponse[NeoliteProDiskState]{}, err
	}

	diskID := neoAliasInt(out, "account_id", "id")
	if diskID == 0 {
		return infer.CreateResponse[NeoliteProDiskState]{},
			fmt.Errorf("create neolite pro disk response tidak ada account_id: %s", neoRawJSON(out))
	}
	id := fmt.Sprintf("%d", diskID)
	partial := NeoliteProDiskState{NeoliteProDiskArgs: in}
	partial.OrderID = neoAliasStr(out, "order_id", "orderid")
	partial.Raw = neoRawJSON(out)

	done, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return c.NeolitePro().DiskGet(ctx, diskID)
		},
		neoliteDiskStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return infer.CreateResponse[NeoliteProDiskState]{ID: id, Output: partial},
				infer.ResourceInitFailedError{Reasons: []string{fmt.Sprintf("neolite pro disk %d belum active: %s", diskID, err)}}
		}
		return infer.CreateResponse[NeoliteProDiskState]{}, err
	}

	state := partial
	state.Status = neoAliasStr(done, "status", "state")
	state.Raw = neoRawJSON(done)
	return infer.CreateResponse[NeoliteProDiskState]{ID: id, Output: state}, nil
}

func (NeoliteProDisk) Update(
	ctx context.Context, req infer.UpdateRequest[NeoliteProDiskArgs, NeoliteProDiskState],
) (infer.UpdateResponse[NeoliteProDiskState], error) {
	if req.DryRun {
		return infer.UpdateResponse[NeoliteProDiskState]{Output: NeoliteProDiskState{NeoliteProDiskArgs: req.Inputs}}, nil
	}

	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.UpdateResponse[NeoliteProDiskState]{}, err
	}
	newSize := i64Val(req.Inputs.Size)
	oldSize := i64Val(req.State.Size)
	if newSize == oldSize {
		return infer.UpdateResponse[NeoliteProDiskState]{Output: req.State}, nil
	}
	if newSize < oldSize {
		return infer.UpdateResponse[NeoliteProDiskState]{},
			fmt.Errorf("neolite pro disk cuma bisa di-upgrade: %d -> %d", oldSize, newSize)
	}

	// upgrade pakai additional_size INCREMENT, bukan target absolute.
	if _, err := c.NeolitePro().DiskUpgrade(ctx, id, client.DiskUpgradePayload{
		AdditionalSize:   newSize - oldSize,
		PayInvoiceWithCC: yesNo(req.Inputs.PayWithCreditCard),
	}); err != nil {
		return infer.UpdateResponse[NeoliteProDiskState]{}, fmt.Errorf("upgrade neolite pro disk: %w", err)
	}

	done, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) {
			return c.NeolitePro().DiskGet(ctx, id)
		},
		neoliteDiskStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return infer.UpdateResponse[NeoliteProDiskState]{
					Output: NeoliteProDiskState{NeoliteProDiskArgs: req.Inputs},
				}, infer.ResourceInitFailedError{Reasons: []string{
					fmt.Sprintf("neolite pro disk %d belum active: %s", id, err),
				}}
		}
		return infer.UpdateResponse[NeoliteProDiskState]{}, fmt.Errorf("neolite pro disk %d gagal balik active: %w", id, err)
	}

	state := NeoliteProDiskState{NeoliteProDiskArgs: req.Inputs}
	state.OrderID = req.State.OrderID
	state.Status = neoAliasStr(done, "status", "state")
	state.Raw = neoRawJSON(done)
	return infer.UpdateResponse[NeoliteProDiskState]{Output: state}, nil
}

func (NeoliteProDisk) Read(
	ctx context.Context, req infer.ReadRequest[NeoliteProDiskArgs, NeoliteProDiskState],
) (infer.ReadResponse[NeoliteProDiskArgs, NeoliteProDiskState], error) {
	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.ReadResponse[NeoliteProDiskArgs, NeoliteProDiskState]{}, err
	}
	out, err := c.NeolitePro().DiskGet(ctx, id)
	if err != nil {
		if client.IsNotFound(err) {
			return infer.ReadResponse[NeoliteProDiskArgs, NeoliteProDiskState]{},
				fmt.Errorf("neolite pro disk %s not found", req.ID)
		}
		return infer.ReadResponse[NeoliteProDiskArgs, NeoliteProDiskState]{}, err
	}

	state := req.State
	state.Status = neoAliasStr(out, "status", "state")
	if v := neoAliasStr(out, "service_name", "name", "label"); v != "" {
		state.ServiceName = &v
	}
	if v := neoAliasInt(out, "size", "disk_size"); v > 0 {
		state.Size = &v
	}
	state.Raw = neoRawJSON(out)
	return infer.ReadResponse[NeoliteProDiskArgs, NeoliteProDiskState]{State: state}, nil
}

func (NeoliteProDisk) Delete(
	ctx context.Context, req infer.DeleteRequest[NeoliteProDiskState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if _, err := c.NeolitePro().DiskDelete(ctx, id); err != nil {
		if client.IsNotFound(err) {
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
