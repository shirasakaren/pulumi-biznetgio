package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type NeoliteSnapshot struct{}

type NeoliteSnapshotArgs struct {
	NeoliteAccountID  int64   `pulumi:"neoliteAccountId"`
	Name              *string `pulumi:"name,optional"`
	Description       *string `pulumi:"description,optional"`
	Cycle             string  `pulumi:"cycle"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
}

type NeoliteSnapshotState struct {
	NeoliteSnapshotArgs
	OrderID string `pulumi:"orderId"`
	Status  string `pulumi:"status"`
}

func (a *NeoliteSnapshotArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.NeoliteAccountID, "Account id VM NEO Lite yang di-snapshot.")
	ann.Describe(&a.Name, "Nama snapshot. Default `snapshot-name`.")
	ann.SetDefault(&a.Name, "snapshot-name")
	ann.Describe(&a.Description, "Deskripsi snapshot.")
	ann.SetDefault(&a.Description, "")
	ann.Describe(&a.Cycle, "Siklus billing: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.PayWithCreditCard, "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). "+
		"Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.Promocode, "Kode promo saat order.")
	ann.SetDefault(&a.Promocode, "")
}

func (s *NeoliteSnapshotState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id dari response create.")
	ann.Describe(&s.Status, "Status snapshot (Active, Pending, Suspended, Terminated).")
}

func (NeoliteSnapshot) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteSnapshotArgs],
) (infer.CreateResponse[NeoliteSnapshotState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteSnapshotState]{
			ID:     "0",
			Output: NeoliteSnapshotState{NeoliteSnapshotArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	billing, err := c.Neolite().SnapshotCreate(ctx, in.NeoliteAccountID, client.SnapshotPayload{
		Cycle:            in.Cycle,
		Name:             strPtr(in.Name),
		Description:      strPtr(in.Description),
		Promocode:        strPtr(in.Promocode),
		PayInvoiceWithCC: yesNo(in.PayWithCreditCard),
	})
	if err != nil {
		return infer.CreateResponse[NeoliteSnapshotState]{}, err
	}
	if billing.AccountID == "" {
		return infer.CreateResponse[NeoliteSnapshotState]{},
			fmt.Errorf("create neolite snapshot response tidak ada account_id: order_id=%s", billing.OrderID)
	}
	id := billing.AccountID
	partial := NeoliteSnapshotState{NeoliteSnapshotArgs: in, OrderID: billing.OrderID}

	acc, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (*client.SnapshotAccountResource, error) {
			n, err := parseNeoID(id)
			if err != nil {
				return nil, err
			}
			return c.Neolite().AccountSnapshotGet(ctx, n)
		},
		neoliteSnapshotStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return infer.CreateResponse[NeoliteSnapshotState]{ID: id, Output: partial},
				infer.ResourceInitFailedError{Reasons: []string{fmt.Sprintf("neolite snapshot %s belum active: %s", id, err)}}
		}
		return infer.CreateResponse[NeoliteSnapshotState]{}, err
	}

	state := partial
	state.Status = acc.Status
	return infer.CreateResponse[NeoliteSnapshotState]{ID: id, Output: state}, nil
}

func (NeoliteSnapshot) Read(
	ctx context.Context, req infer.ReadRequest[NeoliteSnapshotArgs, NeoliteSnapshotState],
) (infer.ReadResponse[NeoliteSnapshotArgs, NeoliteSnapshotState], error) {
	c := GetClient(ctx)
	list, err := c.Neolite().AccountSnapshotList(ctx, "")
	if err != nil {
		return infer.ReadResponse[NeoliteSnapshotArgs, NeoliteSnapshotState]{}, err
	}
	for _, sn := range list {
		if sn.AccountID == req.ID {
			state := req.State
			state.Status = sn.Status
			if sn.ExtraDetails.Name != "" {
				state.Name = &sn.ExtraDetails.Name
			}
			if sn.ExtraDetails.Description != "" {
				state.Description = &sn.ExtraDetails.Description
			}
			return infer.ReadResponse[NeoliteSnapshotArgs, NeoliteSnapshotState]{State: state}, nil
		}
	}
	return infer.ReadResponse[NeoliteSnapshotArgs, NeoliteSnapshotState]{},
		fmt.Errorf("neolite snapshot %s not found", req.ID)
}

func (NeoliteSnapshot) Delete(
	ctx context.Context, req infer.DeleteRequest[NeoliteSnapshotState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if _, err := c.Neolite().SnapshotDelete(ctx, id); err != nil {
		if client.IsNotFound(err) {
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func neoliteSnapshotStatus(a *client.SnapshotAccountResource) string {
	return strings.ToLower(a.Status)
}
