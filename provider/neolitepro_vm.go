package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type NeoliteProVm struct{}

type NeoliteProVmArgs struct {
	SSHAndConsoleUser string  `pulumi:"sshAndConsoleUser"`
	ConsolePassword   string  `pulumi:"consolePassword" provider:"secret"`
	VMName            *string `pulumi:"vmName,optional"`
	Description       *string `pulumi:"description,optional"`
	ProductID         int64   `pulumi:"productId"`
	SelectOS          string  `pulumi:"selectOs"`
	KeypairID         int64   `pulumi:"keypairId"`
	Cycle             string  `pulumi:"cycle"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
	PowerState        *string `pulumi:"powerState,optional"`
	RebuildOS         *string `pulumi:"rebuildOs,optional"`
	DiskSize          *int64  `pulumi:"diskSize,optional"`
}

type NeoliteProVmState struct {
	NeoliteProVmArgs
	OrderID         string             `pulumi:"orderId"`
	Status          string             `pulumi:"status"`
	Uptime          int64              `pulumi:"uptime"`
	MaxDisk         int64              `pulumi:"maxDisk"`
	MaxMem          int64              `pulumi:"maxMem"`
	Mem             int64              `pulumi:"mem"`
	CPUs            int64              `pulumi:"cpus"`
	CIUser          string             `pulumi:"ciUser"`
	CIPassword      string             `pulumi:"ciPassword" provider:"secret"`
	OSName          string             `pulumi:"osName"`
	Region          string             `pulumi:"region"`
	RegionLabel     string             `pulumi:"regionLabel"`
	NextDue         string             `pulumi:"nextDue"`
	RecurringAmount int64              `pulumi:"recurringAmount"`
	Billingcycle    string             `pulumi:"billingcycle"`
	ProductName     string             `pulumi:"productName"`
	LastInvoice     NeoliteLastInvoice `pulumi:"lastInvoice"`
	Raw             string             `pulumi:"raw" provider:"secret"`
}

func (a *NeoliteProVmArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.SSHAndConsoleUser, "User SSH & console yang dipasang saat create.")
	ann.Describe(&a.ConsolePassword, "Password console saat create. Write-only: ga pernah di-refetch dari API.")
	ann.Describe(&a.VMName, "Nama VM. Default `server-name`. Bisa diubah via change-vm-name.")
	ann.SetDefault(&a.VMName, "server-name")
	ann.Describe(&a.Description, "Deskripsi VM.")
	ann.SetDefault(&a.Description, "")
	ann.Describe(&a.ProductID, "Product id dari function `getProducts` atau portal.")
	ann.Describe(&a.SelectOS, "OS yang dipasang saat create, dari function `getOsList`. Ganti OS = pakai `rebuildOs`.")
	ann.Describe(&a.KeypairID, "Id keypair dari `NeoliteProKeypair`. Bisa diganti via change-keypair.")
	ann.Describe(&a.Cycle, "Siklus billing: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.PayWithCreditCard, "Bayar invoice pake kartu kredit saat order. Default true (auto-charge). "+
		"Set false kalau mau ninggalin invoice unpaid di portal — resource bakal stuck Pending sampai dibayar.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.Promocode, "Kode promo saat order.")
	ann.SetDefault(&a.Promocode, "")
	ann.Describe(&a.PowerState, "Power state VM: start, stop, suspend, resume, atau shutdown. "+
		"Update cuma mengirim action kalau nilainya berubah.")
	ann.Describe(&a.RebuildOS, "Kalau berubah, VM di-rebuild (wipe OS) pake OS baru via endpoint rebuild. "+
		"List OS valid ada di function `getOsList`.")
	ann.Describe(&a.DiskSize, "Ukuran disk target (GB, absolute — bukan tambahan). Cuma bisa naik, bukan turun.")
}

func (s *NeoliteProVmState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id dari response create.")
	ann.Describe(&s.Status, "Status akun terakhir dari API (Active, Pending, Suspended, Terminated).")
	ann.Describe(&s.Uptime, "Uptime VM dalam detik.")
	ann.Describe(&s.MaxDisk, "Ukuran disk maksimal VM (GB).")
	ann.Describe(&s.MaxMem, "Memory maksimal VM (MB).")
	ann.Describe(&s.Mem, "Memory yang dipakai VM (MB).")
	ann.Describe(&s.CPUs, "Jumlah CPU VM.")
	ann.Describe(&s.CIUser, "Cloud-init user VM.")
	ann.Describe(&s.CIPassword, "Cloud-init password VM (sensitive).")
	ann.Describe(&s.OSName, "Nama OS yang jalan di VM.")
	ann.Describe(&s.Region, "Region VM.")
	ann.Describe(&s.RegionLabel, "Label region VM.")
	ann.Describe(&s.NextDue, "Tanggal tagihan berikutnya.")
	ann.Describe(&s.RecurringAmount, "Nominal recurring per siklus.")
	ann.Describe(&s.Billingcycle, "Siklus billing aktif.")
	ann.Describe(&s.ProductName, "Nama product aktif.")
	ann.Describe(&s.LastInvoice, "Invoice terakhir VM.")
	ann.Describe(&s.Raw, "Full JSON response akun terakhir dari API, "+
		"buat akses field yang belum dimodel (cipassword di-mask).")
}

func (NeoliteProVm) WireDependencies(f infer.FieldSelector, _ *NeoliteProVmArgs, state *NeoliteProVmState) {
	f.OutputField(&state.ConsolePassword).AlwaysSecret()
	f.OutputField(&state.CIPassword).AlwaysSecret()
}

func (NeoliteProVm) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteProVmArgs],
) (infer.CreateResponse[NeoliteProVmState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteProVmState]{
			ID:     "0",
			Output: NeoliteProVmState{NeoliteProVmArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	billing, err := c.NeolitePro().VMCreate(ctx, client.NeoliteCreatePayload{
		ProductID:         in.ProductID,
		Cycle:             in.Cycle,
		SelectOS:          in.SelectOS,
		KeypairID:         in.KeypairID,
		VMName:            strPtr(in.VMName),
		Description:       strPtr(in.Description),
		SSHAndConsoleUser: in.SSHAndConsoleUser,
		ConsolePassword:   in.ConsolePassword,
		Promocode:         strPtr(in.Promocode),
		PayInvoiceWithCC:  yesNo(in.PayWithCreditCard),
	})
	if err != nil {
		return infer.CreateResponse[NeoliteProVmState]{}, err
