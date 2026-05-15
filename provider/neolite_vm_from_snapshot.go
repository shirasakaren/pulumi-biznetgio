package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type NeoliteVmFromSnapshot struct{}

type NeoliteVmFromSnapshotArgs struct {
	SnapshotID        int64   `pulumi:"snapshotId"`
	ProductID         int64   `pulumi:"productId"`
	Cycle             string  `pulumi:"cycle"`
	KeypairID         int64   `pulumi:"keypairId"`
	Name              string  `pulumi:"name"`
	Description       *string `pulumi:"description,optional"`
	SSHAndConsoleUser string  `pulumi:"sshAndConsoleUser"`
	ConsolePassword   string  `pulumi:"consolePassword" provider:"secret"`
	Promocode         *string `pulumi:"promocode,optional"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
}

type NeoliteVmFromSnapshotState struct {
	NeoliteVmFromSnapshotArgs
	OrderID string `pulumi:"orderId"`
	Status  string `pulumi:"status"`
}

func (a *NeoliteVmFromSnapshotArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.SnapshotID, "Account id snapshot sumber, dari `NeoliteSnapshot`.")
	ann.Describe(&a.ProductID, "Product id dari function `getProducts` atau portal.")
	ann.Describe(&a.Cycle, "Siklus billing: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.KeypairID, "Id keypair dari `NeoliteKeypair`.")
	ann.Describe(&a.Name, "Nama VM hasil restore.")
	ann.Describe(&a.Description, "Deskripsi VM.")
	ann.SetDefault(&a.Description, "")
	ann.Describe(&a.SSHAndConsoleUser, "User SSH & console yang dipasang saat create.")
	ann.Describe(&a.ConsolePassword, "Password console saat create. Write-only: ga pernah di-refetch dari API.")
	ann.Describe(&a.Promocode, "Kode promo saat order.")
	ann.SetDefault(&a.Promocode, "")
	ann.Describe(&a.PayWithCreditCard, "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). "+
		"Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.")
	ann.SetDefault(&a.PayWithCreditCard, true)
}

func (s *NeoliteVmFromSnapshotState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id dari response create.")
	ann.Describe(&s.Status, "Status VM (Active, Pending, Suspended, Terminated).")
}

func (NeoliteVmFromSnapshot) WireDependencies(
	f infer.FieldSelector, _ *NeoliteVmFromSnapshotArgs, state *NeoliteVmFromSnapshotState,
) {
	f.OutputField(&state.ConsolePassword).AlwaysSecret()
}

func (NeoliteVmFromSnapshot) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteVmFromSnapshotArgs],
) (infer.CreateResponse[NeoliteVmFromSnapshotState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteVmFromSnapshotState]{
			ID:     "0",
			Output: NeoliteVmFromSnapshotState{NeoliteVmFromSnapshotArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	billing, err := c.Neolite().SnapshotRestoreWith(ctx, in.SnapshotID, client.NeoliteFromSnapshotPayload{
		ProductID:         in.ProductID,
		Cycle:             in.Cycle,
		KeypairID:         in.KeypairID,
		Name:              in.Name,
		Description:       strPtr(in.Description),
