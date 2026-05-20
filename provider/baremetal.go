package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type Baremetal struct{}

type BaremetalArgs struct {
	ProductID         int64   `pulumi:"productId" provider:"replaceOnChanges"`
	Cycle             string  `pulumi:"cycle" provider:"replaceOnChanges"`
	SelectOS          *string `pulumi:"selectOs,optional" provider:"replaceOnChanges"`
	KeypairID         int64   `pulumi:"keypairId" provider:"replaceOnChanges"`
	Label             string  `pulumi:"label"`
	PublicIP          *int64  `pulumi:"publicIp,optional" provider:"replaceOnChanges"`
	Promocode         *string `pulumi:"promocode,optional" provider:"replaceOnChanges"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional" provider:"replaceOnChanges"`
	PowerState        *string `pulumi:"powerState,optional"`
	ResetTrigger      *string `pulumi:"resetTrigger,optional"`
	RebuildOS         *string `pulumi:"rebuildOs,optional"`
}

type BaremetalState struct {
	BaremetalArgs
	Status    string  `pulumi:"status"`
	OrderID   *string `pulumi:"orderId"`
	IPAddress *string `pulumi:"ipAddress"`
	CreatedAt *string `pulumi:"createdAt"`
	Raw       string  `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id from the baremetal products function.")
	ann.Describe(&a.Cycle, "Billing cycle, e.g. `m` for monthly or `a` for annual.")
	ann.Describe(&a.SelectOS, "OS to install at creation, from `GET /baremetals/products/{product_id}/oss`. "+
		"Defaults to `ubuntu-22`. Create-only, changing it replaces the instance.")
	ann.SetDefault(&a.SelectOS, "ubuntu-22")
	ann.Describe(&a.KeypairID, "Baremetal keypair id, from BaremetalKeypair. "+
		"The baremetal keypair pool is separate from neolite/gpu.")
	ann.Describe(&a.Label, "Display name of the server. The only field updatable in place.")
	ann.Describe(&a.PublicIP, "Number of public IPs to request (1 = with public IP). Defaults to 1. "+
		"Create-only, changing it replaces the instance.")
	ann.SetDefault(&a.PublicIP, int64(1))
	ann.Describe(&a.Promocode, "Promo code to apply at creation.")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card. Defaults to true. "+
		"Set false to leave the invoice unpaid in the portal; the resource stays Pending until paid.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.PowerState, "Power state of the server: `on` or `off`. Only fires an API call when the value changes.")
	ann.Describe(&a.ResetTrigger, "Change this value to trigger a one-shot reset/reboot. "+
		"The reset is an action, not a stable state.")
	ann.Describe(&a.RebuildOS, "Change this value to rebuild the instance (destructive, wipes the OS) "+
		"via `PUT /baremetals/{account_id}/rebuild`. Valid OS list comes from the rebuildOsList function.")
}

func (s *BaremetalState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Status, "Current status of the server.")
	ann.Describe(&s.OrderID, "Order id from the creation response.")
	ann.Describe(&s.IPAddress, "Public IP address of the server, when present.")
	ann.Describe(&s.CreatedAt, "Creation timestamp (alias of created_at/date_created).")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (Baremetal) Create(
	ctx context.Context, req infer.CreateRequest[BaremetalArgs],
) (infer.CreateResponse[BaremetalState], error) {
	resp := infer.CreateResponse[BaremetalState]{Output: baremetalStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
