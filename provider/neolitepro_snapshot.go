package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type NeoliteProSnapshot struct{}

type NeoliteProSnapshotArgs struct {
	NeoliteAccountID  int64   `pulumi:"neoliteAccountId"`
	Name              *string `pulumi:"name,optional"`
	Description       *string `pulumi:"description,optional"`
	Cycle             string  `pulumi:"cycle"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
}

type NeoliteProSnapshotState struct {
	NeoliteProSnapshotArgs
	OrderID string `pulumi:"orderId"`
	Status  string `pulumi:"status"`
}

func (a *NeoliteProSnapshotArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.NeoliteAccountID, "Account id VM NEO Lite Pro yang di-snapshot.")
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

func (s *NeoliteProSnapshotState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id dari response create.")
	ann.Describe(&s.Status, "Status snapshot (Active, Pending, Suspended, Terminated).")
}

func (NeoliteProSnapshot) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteProSnapshotArgs],
) (infer.CreateResponse[NeoliteProSnapshotState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteProSnapshotState]{
			ID:     "0",
			Output: NeoliteProSnapshotState{NeoliteProSnapshotArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	billing, err := c.NeolitePro().SnapshotCreate(ctx, in.NeoliteAccountID, client.SnapshotPayload{
		Cycle:            in.Cycle,
		Name:             strPtr(in.Name),
		Description:      strPtr(in.Description),
		Promocode:        strPtr(in.Promocode),
		PayInvoiceWithCC: yesNo(in.PayWithCreditCard),
	})
	if err != nil {
		return infer.CreateResponse[NeoliteProSnapshotState]{}, err
	}
	if billing.AccountID == "" {
