package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type GpuInstance struct{}

type GpuSubscriptionArgs struct {
	Cycle string `pulumi:"cycle"`
}

type GpuOnDemandArgs struct {
	AdditionalHours *int64 `pulumi:"additionalHours,optional"`
}

type GpuInstanceArgs struct {
	ProductID                     int64                `pulumi:"productId" provider:"replaceOnChanges"`
	SelectOS                      string               `pulumi:"selectOs" provider:"replaceOnChanges"`
	KeypairID                     int64                `pulumi:"keypairId" provider:"replaceOnChanges"`
	ServiceName                   *string              `pulumi:"serviceName,optional" provider:"replaceOnChanges"`
	SSHAndConsoleUser             string               `pulumi:"sshAndConsoleUser" provider:"replaceOnChanges"`
	ConsolePassword               string               `pulumi:"consolePassword" provider:"secret,replaceOnChanges"` //nolint:lll
	Promocode                     *string              `pulumi:"promocode,optional" provider:"replaceOnChanges"`
	PayWithCreditCard             *bool                `pulumi:"payWithCreditCard,optional" provider:"replaceOnChanges"`
	Subscription                  *GpuSubscriptionArgs `pulumi:"subscription,optional" provider:"replaceOnChanges"`
	OnDemand                      *GpuOnDemandArgs     `pulumi:"onDemand,optional" provider:"replaceOnChanges"`
	RebuildTrigger                *string              `pulumi:"rebuildTrigger,optional"`
	ReserveAdditionalHoursTrigger *string              `pulumi:"reserveAdditionalHoursTrigger,optional"`
}

type GpuInstanceState struct {
	GpuInstanceArgs
	Status  string  `pulumi:"status"`
	OrderID *string `pulumi:"orderId"`
	Raw     string  `pulumi:"raw" provider:"secret"`
}

func (a *GpuSubscriptionArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.Cycle, "Billing cycle, e.g. `m` for monthly or `a` for annual.")
}

func (a *GpuOnDemandArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AdditionalHours, "Additional hours to reserve. Defaults to 0.")
}

func (a *GpuInstanceArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "GPU product id from the gpu products function.")
	ann.Describe(&a.SelectOS, "OS to install, from the product's select-os catalog.")
	ann.Describe(&a.KeypairID, "GPU keypair id, from GpuKeypair.")
	ann.Describe(&a.ServiceName, "Display name of the instance. Create-only, changing it replaces the instance.")
	ann.Describe(&a.SSHAndConsoleUser, "SSH and console user for the instance.")
	ann.Describe(&a.ConsolePassword, "Console password. Create-only, changing it replaces the instance.")
	ann.Describe(&a.Promocode, "Promo code to apply at creation.")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card. Defaults to true. "+
		"Set false to leave the invoice unpaid in the portal; the resource stays Pending until paid.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.Subscription, "Subscription billing mode (cycle-based). "+
		"Exactly one of `subscription` or `onDemand` must be set.")
	ann.Describe(&a.OnDemand, "On-demand billing mode (hourly). Exactly one of `subscription` or `onDemand` must be set.")
	ann.Describe(&a.RebuildTrigger, "Change this value to rebuild the instance (destructive, reinstalls the OS).")
	ann.Describe(&a.ReserveAdditionalHoursTrigger, "Change this value to reserve 1 additional hour "+
		"on an on-demand instance.")
}

func (s *GpuInstanceState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Status, "Current status of the instance.")
	ann.Describe(&s.OrderID, "Order id from the creation response.")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (GpuInstance) WireDependencies(f infer.FieldSelector, _ *GpuInstanceArgs, state *GpuInstanceState) {
	f.OutputField(&state.ConsolePassword).AlwaysSecret()
}

func (GpuInstance) Check(ctx context.Context, req infer.CheckRequest) (infer.CheckResponse[GpuInstanceArgs], error) {
	inputs, failures, err := infer.DefaultCheck[GpuInstanceArgs](ctx, req.NewInputs)
	if err != nil {
		return infer.CheckResponse[GpuInstanceArgs]{}, err
