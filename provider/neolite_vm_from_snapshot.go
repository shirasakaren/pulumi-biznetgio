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
