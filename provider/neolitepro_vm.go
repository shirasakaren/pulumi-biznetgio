package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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
	ann.Describe(&a.SSHAndConsoleUser, "SSH and console user set at creation.")
	ann.Describe(&a.ConsolePassword, "Console password at creation. Write-only: never re-fetched from the API.")
	ann.Describe(&a.VMName, "VM name. Defaults to `server-name`. Can be changed via change-vm-name.")
	ann.SetDefault(&a.VMName, "server-name")
	ann.Describe(&a.Description, "VM description.")
	ann.SetDefault(&a.Description, "")
	ann.Describe(&a.ProductID, "Product id from the `getProducts` function or the portal.")
	ann.Describe(&a.SelectOS, "OS installed at creation, from the `getOsList` function. To change OS, use `rebuildOs`.")
	ann.Describe(&a.KeypairID, "Keypair id from `NeoliteProKeypair`. Can be changed via change-keypair.")
	ann.Describe(&a.Cycle, "Billing cycle: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card at order time. "+
		"Defaults to true (auto-charge); set false to leave it unpaid in the portal until settled.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.Promocode, "Promo code to apply at order.")
	ann.SetDefault(&a.Promocode, "")
	ann.Describe(&a.PowerState, "VM power state: start, stop, suspend, resume, or shutdown. "+
		"Update only sends an action when the value changes.")
	ann.Describe(&a.RebuildOS, "When changed, the VM is rebuilt (wipes the OS) with the new OS via the rebuild endpoint. "+
		"Valid OS values are listed by the `getOsList` function.")
	ann.Describe(&a.DiskSize, "Target disk size (GB, absolute - not an increment). Can only go up, never down.")
}

func (s *NeoliteProVmState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id from the creation response.")
	ann.Describe(&s.Status, "Last known account status from the API (Active, Pending, Suspended, Terminated).")
	ann.Describe(&s.Uptime, "VM uptime in seconds.")
	ann.Describe(&s.MaxDisk, "Maximum VM disk size (GB).")
	ann.Describe(&s.MaxMem, "Maximum VM memory (MB).")
	ann.Describe(&s.Mem, "Memory used by the VM (MB).")
	ann.Describe(&s.CPUs, "Number of VM CPUs.")
	ann.Describe(&s.CIUser, "VM cloud-init user.")
	ann.Describe(&s.CIPassword, "VM cloud-init password (sensitive).")
	ann.Describe(&s.OSName, "Name of the OS running on the VM.")
	ann.Describe(&s.Region, "VM region.")
	ann.Describe(&s.RegionLabel, "VM region label.")
	ann.Describe(&s.NextDue, "Next billing due date.")
	ann.Describe(&s.RecurringAmount, "Recurring amount per cycle.")
	ann.Describe(&s.Billingcycle, "Active billing cycle.")
	ann.Describe(&s.ProductName, "Active product name.")
	ann.Describe(&s.LastInvoice, "VM's last invoice.")
	ann.Describe(&s.Raw, "Raw JSON response of the last-read account, "+
		"for accessing fields not yet modeled (cipassword masked).")
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
	}
	if billing.AccountID == "" {
		return infer.CreateResponse[NeoliteProVmState]{},
			fmt.Errorf("create neolite pro vm response tidak ada account_id: order_id=%s", billing.OrderID)
	}
	id := billing.AccountID
	partial := NeoliteProVmState{NeoliteProVmArgs: in, OrderID: billing.OrderID}

	_, err = client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (*client.AccountResource, error) {
			n, err := parseNeoID(id)
			if err != nil {
				return nil, err
			}
			return c.NeolitePro().AccountGet(ctx, n)
		},
		neoliteStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return infer.CreateResponse[NeoliteProVmState]{ID: id, Output: partial},
				infer.ResourceInitFailedError{Reasons: []string{fmt.Sprintf("neolite pro vm %s not active yet: %s", id, err)}}
		}
		return infer.CreateResponse[NeoliteProVmState]{}, err
	}

	state, err := readNeoliteProVm(ctx, c, id, in, partial)
	if err != nil {
		return infer.CreateResponse[NeoliteProVmState]{}, err
	}
	return infer.CreateResponse[NeoliteProVmState]{ID: id, Output: state}, nil
}

func (NeoliteProVm) Update(
	ctx context.Context, req infer.UpdateRequest[NeoliteProVmArgs, NeoliteProVmState],
) (infer.UpdateResponse[NeoliteProVmState], error) {
	if req.DryRun {
		return infer.UpdateResponse[NeoliteProVmState]{Output: NeoliteProVmState{NeoliteProVmArgs: req.Inputs}}, nil
	}

	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.UpdateResponse[NeoliteProVmState]{}, err
	}
	in := req.Inputs
	old := req.State

	if !eqStrPtr(in.VMName, old.VMName) {
		if _, err := c.NeolitePro().VMChangeName(ctx, id, strPtr(in.VMName)); err != nil {
			return infer.UpdateResponse[NeoliteProVmState]{}, fmt.Errorf("change neolite pro vm name: %w", err)
		}
	}
	if in.KeypairID != old.KeypairID {
		if _, err := c.NeolitePro().VMChangeKeypair(ctx, id, in.KeypairID); err != nil {
			return infer.UpdateResponse[NeoliteProVmState]{}, fmt.Errorf("change neolite pro vm keypair: %w", err)
		}
	}

	needsPoll := false
	if in.ProductID != old.ProductID {
		if _, err := c.NeolitePro().VMChangePackage(ctx, id, client.ChangePackagePayload{
			NewProductID:     in.ProductID,
			PayInvoiceWithCC: yesNo(in.PayWithCreditCard),
		}); err != nil {
			return infer.UpdateResponse[NeoliteProVmState]{}, fmt.Errorf("change neolite pro vm package: %w", err)
		}
		needsPoll = true
	}
	if !eqI64Ptr(in.DiskSize, old.DiskSize) {
		newSize := i64Val(in.DiskSize)
		oldSize := i64Val(old.DiskSize)
		if newSize < oldSize {
			return infer.UpdateResponse[NeoliteProVmState]{},
				fmt.Errorf("neolite pro vm storage can only be upgraded: %d -> %d", oldSize, newSize)
		}
		if _, err := c.NeolitePro().VMChangeStorage(ctx, id, client.UpgradePayload{
			DiskSize:         newSize,
			PayInvoiceWithCC: yesNo(in.PayWithCreditCard),
		}); err != nil {
			return infer.UpdateResponse[NeoliteProVmState]{}, fmt.Errorf("change neolite pro vm storage: %w", err)
		}
		needsPoll = true
	}
	if !eqStrPtr(in.PowerState, old.PowerState) {
		if _, err := c.NeolitePro().VMState(ctx, id, strPtr(in.PowerState)); err != nil {
			return infer.UpdateResponse[NeoliteProVmState]{},
				fmt.Errorf("set neolite pro vm power state %q: %w", strPtr(in.PowerState), err)
		}
	}
	if !eqStrPtr(in.RebuildOS, old.RebuildOS) {
		if _, err := c.NeolitePro().VMRebuild(ctx, id, strPtr(in.RebuildOS)); err != nil {
			return infer.UpdateResponse[NeoliteProVmState]{}, fmt.Errorf("rebuild neolite pro vm: %w", err)
		}
		needsPoll = true
	}

	if needsPoll {
		_, err := client.WaitForStatus(ctx, 5*time.Second,
			func(ctx context.Context) (*client.AccountResource, error) {
				return c.NeolitePro().AccountGet(ctx, id)
			},
			neoliteStatus, []string{"active"}, []string{"suspended", "terminated"})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return infer.UpdateResponse[NeoliteProVmState]{
						Output: NeoliteProVmState{NeoliteProVmArgs: in},
					}, infer.ResourceInitFailedError{Reasons: []string{
						fmt.Sprintf("neolite pro vm %d not active yet: %s", id, err),
					}}
			}
			return infer.UpdateResponse[NeoliteProVmState]{}, fmt.Errorf("neolite pro vm %d gagal balik active: %w", id, err)
		}
	}

	state, err := readNeoliteProVm(ctx, c, req.ID, in, old)
	if err != nil {
		return infer.UpdateResponse[NeoliteProVmState]{}, err
	}
	return infer.UpdateResponse[NeoliteProVmState]{Output: state}, nil
}

func (NeoliteProVm) Read(
	ctx context.Context, req infer.ReadRequest[NeoliteProVmArgs, NeoliteProVmState],
) (infer.ReadResponse[NeoliteProVmArgs, NeoliteProVmState], error) {
	c := GetClient(ctx)
	state, err := readNeoliteProVm(ctx, c, req.ID, req.State.NeoliteProVmArgs, req.State)
	if err != nil {
		return infer.ReadResponse[NeoliteProVmArgs, NeoliteProVmState]{}, err
	}
	return infer.ReadResponse[NeoliteProVmArgs, NeoliteProVmState]{State: state}, nil
}

func (NeoliteProVm) Delete(
	ctx context.Context, req infer.DeleteRequest[NeoliteProVmState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if _, err := c.NeolitePro().VMDelete(ctx, id); err != nil {
		if client.IsNotFound(err) {
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func readNeoliteProVm(
	ctx context.Context, c *client.Client, id string, in NeoliteProVmArgs, prev NeoliteProVmState,
) (NeoliteProVmState, error) {
	n, err := parseNeoID(id)
	if err != nil {
		return prev, fmt.Errorf("neolite pro vm %s invalid id: %w", id, err)
	}
	acc, err := c.NeolitePro().AccountGet(ctx, n)
	if err != nil {
		if client.IsNotFound(err) {
			return prev, fmt.Errorf("neolite pro vm %s not found", id)
		}
		return prev, err
	}

	st := NeoliteProVmState{NeoliteProVmArgs: in}
	st.OrderID = prev.OrderID
	st.Status = acc.Status
	st.Billingcycle = acc.Billingcycle
	st.NextDue = acc.NextDue
	st.RecurringAmount = acc.RecurringAmount
	st.ProductID = acc.ProductID
	st.ProductName = acc.ProductName
	st.Description = &acc.Description
	st.Region = acc.ExtraDetails.Region
	st.RegionLabel = acc.ExtraDetails.RegionLabel
	st.CIUser = acc.ExtraDetails.CIUser
	st.CIPassword = acc.ExtraDetails.CIPassword
	st.OSName = acc.ExtraDetails.OSName
	if acc.ExtraDetails.KeypairID != 0 {
		st.KeypairID = acc.ExtraDetails.KeypairID
	}
	if v := acc.ExtraDetails.DiskSize; v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			st.DiskSize = &n
		}
	}
	if v := acc.ExtraDetails.Name; v != "" {
		st.VMName = &v
	}
	st.LastInvoice = NeoliteLastInvoice{
		ID:          acc.LastInvoice.ID,
		PaidID:      acc.LastInvoice.PaidID,
		Status:      acc.LastInvoice.Status,
		Date:        acc.LastInvoice.Date,
		Duedate:     acc.LastInvoice.Duedate,
		Paybefore:   acc.LastInvoice.Paybefore,
		Datepaid:    acc.LastInvoice.Datepaid,
		InvoiceType: acc.LastInvoice.InvoiceType,
	}

	if n, err := parseNeoID(id); err == nil {
		if vm, err := c.NeolitePro().VMDetails(ctx, n); err == nil {
			st.Uptime = vm.Uptime
			st.MaxDisk = vm.MaxDisk
			st.MaxMem = vm.MaxMem
			st.Mem = vm.Mem
			st.CPUs = vm.CPUs
		}
	}

	if b, err := json.Marshal(acc); err == nil {
		var m map[string]any
		if json.Unmarshal(b, &m) == nil {
			st.Raw = string(RedactJSON(m))
		}
	}
	return st, nil
}
