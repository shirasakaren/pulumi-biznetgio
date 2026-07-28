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
	}
	if (inputs.Subscription == nil) == (inputs.OnDemand == nil) {
		failures = append(failures, p.CheckFailure{
			Property: "subscription",
			Reason:   "exactly one of subscription or onDemand must be set",
		})
	}
	return infer.CheckResponse[GpuInstanceArgs]{Inputs: inputs, Failures: failures}, nil
}

func (GpuInstance) Create(
	ctx context.Context, req infer.CreateRequest[GpuInstanceArgs],
) (infer.CreateResponse[GpuInstanceState], error) {
	resp := infer.CreateResponse[GpuInstanceState]{Output: gpuStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs

	raw, err := gpuCreate(ctx, c, a)
	if err != nil {
		return infer.CreateResponse[GpuInstanceState]{}, err
	}
	accountID, ok := gpuInt64(raw, "account_id", "id")
	if !ok {
		return infer.CreateResponse[GpuInstanceState]{},
			fmt.Errorf("biznetgio: gpu create response missing account_id: %s", gpuJSON(raw))
	}
	resp.ID = strconv.FormatInt(accountID, 10)
	resp.Output = gpuStateFromMap(ctx, a, raw)

	if _, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) { return c.GPU().AccountGet(ctx, accountID) },
		gpuStatus, []string{"active"}, []string{"terminated", "failed", "error", "deleted", "cancelled"}); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return resp, infer.ResourceInitFailedError{Reasons: []string{
				fmt.Sprintf("gpu instance %d belum active, lanjutin via update aja: %s", accountID, err),
			}}
		}
		return resp, err
	}
	return resp, nil
}

func (GpuInstance) Update(
	ctx context.Context, req infer.UpdateRequest[GpuInstanceArgs, GpuInstanceState],
) (infer.UpdateResponse[GpuInstanceState], error) {
	resp := infer.UpdateResponse[GpuInstanceState]{Output: gpuStateFromMap(ctx, req.Inputs, nil)}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.UpdateResponse[GpuInstanceState]{}, fmt.Errorf("biznetgio: invalid gpu instance id %q: %s", req.ID, err)
	}

	if t := gpuStr(a.RebuildTrigger); t != "" && t != gpuStr(req.State.RebuildTrigger) {
		if _, err := c.GPU().Rebuild(ctx, accountID, a.SelectOS); err != nil {
			return infer.UpdateResponse[GpuInstanceState]{}, fmt.Errorf("biznetgio: rebuild gpu instance %d: %w", accountID, err)
		}
	}
	if t := gpuStr(a.ReserveAdditionalHoursTrigger); t != "" && t != gpuStr(req.State.ReserveAdditionalHoursTrigger) {
		hours := int64(1)
		if a.OnDemand != nil && a.OnDemand.AdditionalHours != nil && *a.OnDemand.AdditionalHours > 0 {
			hours = *a.OnDemand.AdditionalHours
		}
		if _, err := c.GPU().ReserveAdditionalHours(ctx, accountID, hours); err != nil {
			return infer.UpdateResponse[GpuInstanceState]{},
				fmt.Errorf("biznetgio: reserve additional hours for gpu instance %d: %w", accountID, err)
		}
	}

	final, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (map[string]any, error) { return c.GPU().AccountGet(ctx, accountID) },
		gpuStatus, []string{"active"}, []string{"terminated", "failed", "error", "deleted", "cancelled"})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return infer.UpdateResponse[GpuInstanceState]{
				Output: gpuStateFromMap(ctx, a, nil),
			}, infer.ResourceInitFailedError{Reasons: []string{err.Error()}}
		}
		return infer.UpdateResponse[GpuInstanceState]{}, err
	}
	resp.Output = gpuStateFromMap(ctx, a, final)
	return resp, nil
}

