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
