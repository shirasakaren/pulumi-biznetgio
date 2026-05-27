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
	ann.Describe(&a.PayWithCreditCard, "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). "+
		"Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.Size, "Ukuran disk (GB). Default 60. Cuma bisa naik, bukan turun.")
	ann.SetDefault(&a.Size, int64(60))
}

func (s *NeoliteDiskState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id dari response create.")
	ann.Describe(&s.Status, "Status disk (Active, Pending, Suspended, Terminated).")
	ann.Describe(&s.Raw, "Full JSON response disk terakhir dari API, buat akses field yang belum dimodel.")
}

func (NeoliteDisk) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteDiskArgs],
) (infer.CreateResponse[NeoliteDiskState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteDiskState]{
			ID:     "0",
			Output: NeoliteDiskState{NeoliteDiskArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	out, err := c.Neolite().DiskCreate(ctx, client.NeoliteDiskCreatePayload{
		ProductID:        in.ProductID,
		Cycle:            in.Cycle,
		NeoliteAccountID: in.NeoliteAccountID,
		ServiceName:      strPtr(in.ServiceName),
		Promocode:        strPtr(in.Promocode),
		PayInvoiceWithCC: yesNo(in.PayWithCreditCard),
		Size:             i64Val(in.Size),
	})
	if err != nil {
		return infer.CreateResponse[NeoliteDiskState]{}, err
	}

	diskID := neoAliasInt(out, "account_id", "id")
	if diskID == 0 {
		return infer.CreateResponse[NeoliteDiskState]{},
			fmt.Errorf("create neolite disk response tidak ada account_id: %s", neoRawJSON(out))
