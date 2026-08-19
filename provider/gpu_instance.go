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

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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

func (GpuInstance) Read(
	ctx context.Context, req infer.ReadRequest[GpuInstanceArgs, GpuInstanceState],
) (infer.ReadResponse[GpuInstanceArgs, GpuInstanceState], error) {
	resp := infer.ReadResponse[GpuInstanceArgs, GpuInstanceState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  gpuStateFromMap(ctx, req.Inputs, nil),
	}
	c := GetClient(ctx)
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return resp, fmt.Errorf("biznetgio: invalid gpu instance id %q: %s", req.ID, err)
	}
	m, err := c.GPU().AccountGet(ctx, accountID)
	if err != nil {
		if client.IsNotFound(err) {
			return resp, fmt.Errorf("biznetgio: gpu instance %s not found", req.ID)
		}
		return resp, err
	}
	resp.State = gpuStateFromMap(ctx, req.Inputs, m)
	return resp, nil
}

func (GpuInstance) Delete(
	ctx context.Context, req infer.DeleteRequest[GpuInstanceState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	accountID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return infer.DeleteResponse{}, fmt.Errorf("biznetgio: invalid gpu instance id %q: %s", req.ID, err)
	}
	if _, err := c.GPU().Delete(ctx, accountID); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func gpuCreate(ctx context.Context, c *client.Client, a GpuInstanceArgs) (map[string]any, error) {
	pay := "yes"
	if a.PayWithCreditCard != nil && !*a.PayWithCreditCard {
		pay = "no"
	}
	if a.Subscription != nil {
		return c.GPU().Create(ctx, client.NEOGPUCreatePayload{
			ProductID:         a.ProductID,
			SelectOS:          a.SelectOS,
			KeypairID:         a.KeypairID,
			ServiceName:       gpuStr(a.ServiceName),
			SSHAndConsoleUser: a.SSHAndConsoleUser,
			ConsolePassword:   a.ConsolePassword,
			Promocode:         gpuStr(a.Promocode),
			PayInvoiceWithCC:  pay,
			Cycle:             a.Subscription.Cycle,
		})
	}
	if a.OnDemand != nil {
		hours := int64(0)
		if a.OnDemand.AdditionalHours != nil {
			hours = *a.OnDemand.AdditionalHours
		}
		return c.GPU().CreateOneTime(ctx, client.NEOGPUOneTimeCreatePayload{
			ProductID:         a.ProductID,
			SelectOS:          a.SelectOS,
			KeypairID:         a.KeypairID,
			ServiceName:       gpuStr(a.ServiceName),
			SSHAndConsoleUser: a.SSHAndConsoleUser,
			ConsolePassword:   a.ConsolePassword,
			Promocode:         gpuStr(a.Promocode),
			PayInvoiceWithCC:  pay,
			AdditionalHours:   hours,
		})
	}
	return nil, fmt.Errorf("exactly one of subscription or onDemand must be set")
}

func gpuStateFromMap(_ context.Context, args GpuInstanceArgs, m map[string]any) GpuInstanceState {
	st := GpuInstanceState{GpuInstanceArgs: args}
	if m == nil {
		return st
	}
	if v, ok := gpuString(m, "service_name", "name", "label"); ok {
		st.ServiceName = &v
	}
	if v, ok := gpuInt64(m, "product_id"); ok {
		st.ProductID = v
	}
	if v, ok := gpuInt64(m, "keypair_id", "neosshkey_id"); ok {
		st.KeypairID = v
	}
	if v, ok := gpuString(m, "select_os", "os"); ok {
		st.SelectOS = v
	}
	if v, ok := gpuString(m, "ssh_and_console_user", "ciuser", "user"); ok {
		st.SSHAndConsoleUser = v
	}
	if v, ok := gpuString(m, "console_password", "cipassword"); ok {
		st.ConsolePassword = v
	}
	if v, ok := gpuString(m, "cycle", "billingcycle"); ok {
		st.Subscription = &GpuSubscriptionArgs{Cycle: v}
	}
	if v, ok := gpuInt64(m, "additional_hours"); ok {
		st.OnDemand = &GpuOnDemandArgs{AdditionalHours: &v}
	}
	if v, ok := gpuString(m, "order_id"); ok {
		st.OrderID = &v
	}
	st.Status = gpuStringDefault(m, "status", "state")
	st.Raw = string(gpuJSON(m))
	return st
}

func gpuStatus(m map[string]any) string {
	return strings.ToLower(gpuStringDefault(m, "status", "state"))
}

func gpuStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func gpuInt64(m map[string]any, keys ...string) (int64, bool) {
	for _, k := range keys {
		v, ok := m[k]
		if !ok {
			continue
		}
		switch n := v.(type) {
		case float64:
			return int64(n), true
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return i, true
			}
		case string:
			if i, err := strconv.ParseInt(n, 10, 64); err == nil {
				return i, true
			}
		}
	}
	return 0, false
}

func gpuString(m map[string]any, keys ...string) (string, bool) {
	for _, k := range keys {
		v, ok := m[k]
		if !ok {
			continue
		}
		switch s := v.(type) {
		case string:
			return s, true
		case json.Number:
			return s.String(), true
		case float64:
			return strconv.FormatFloat(s, 'f', -1, 64), true
		case bool:
			return strconv.FormatBool(s), true
		}
	}
	return "", false
}

func gpuStringDefault(m map[string]any, keys ...string) string {
	s, _ := gpuString(m, keys...)
	return s
}

func gpuJSON(m map[string]any) []byte {
	return RedactJSON(m)
}
