package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type NeoliteDisk struct{}

type NeoliteDiskArgs struct {
	ProductID         int64   `pulumi:"productId"`
	Cycle             string  `pulumi:"cycle"`
	NeoliteAccountID  int64   `pulumi:"neoliteAccountId"`
	ServiceName       *string `pulumi:"serviceName,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
	Size              *int64  `pulumi:"size,optional"`
}

type NeoliteDiskState struct {
	NeoliteDiskArgs
	OrderID string `pulumi:"orderId"`
	Status  string `pulumi:"status"`
	Raw     string `pulumi:"raw" provider:"secret"`
}

func (a *NeoliteDiskArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id disk dari endpoint `/neolites/disks/products`.")
	ann.Describe(&a.Cycle, "Siklus billing: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.NeoliteAccountID, "Account id VM NEO Lite tempat disk dipasang.")
	ann.Describe(&a.ServiceName, "Nama layanan disk. Default `service-name`.")
	ann.SetDefault(&a.ServiceName, "service-name")
	ann.Describe(&a.Promocode, "Kode promo saat order.")
	ann.SetDefault(&a.Promocode, "")
