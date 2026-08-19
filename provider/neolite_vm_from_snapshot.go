package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
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
	ann.Describe(&a.SnapshotID, "Account id of the source snapshot, from `NeoliteSnapshot`.")
	ann.Describe(&a.ProductID, "Product id from the `getProducts` function or the portal.")
	ann.Describe(&a.Cycle, "Billing cycle: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.KeypairID, "Keypair id from `NeoliteKeypair`.")
	ann.Describe(&a.Name, "Name of the restored VM.")
	ann.Describe(&a.Description, "VM description.")
	ann.SetDefault(&a.Description, "")
	ann.Describe(&a.SSHAndConsoleUser, "SSH and console user set at creation.")
	ann.Describe(&a.ConsolePassword, "Console password at creation. Write-only: never re-fetched from the API.")
	ann.Describe(&a.Promocode, "Promo code to apply at order.")
	ann.SetDefault(&a.Promocode, "")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card at order time. "+
		"Defaults to true (auto-charge); set false to leave it unpaid in the portal until settled.")
	ann.SetDefault(&a.PayWithCreditCard, true)
}

func (s *NeoliteVmFromSnapshotState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id from the creation response.")
	ann.Describe(&s.Status, "VM status (Active, Pending, Suspended, Terminated).")
}

func (NeoliteVmFromSnapshot) WireDependencies(
	f infer.FieldSelector, _ *NeoliteVmFromSnapshotArgs, state *NeoliteVmFromSnapshotState,
) {
	f.OutputField(&state.ConsolePassword).AlwaysSecret()
}

func (NeoliteVmFromSnapshot) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteVmFromSnapshotArgs],
) (infer.CreateResponse[NeoliteVmFromSnapshotState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteVmFromSnapshotState]{
			ID:     "0",
			Output: NeoliteVmFromSnapshotState{NeoliteVmFromSnapshotArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	billing, err := c.Neolite().SnapshotRestoreWith(ctx, in.SnapshotID, client.NeoliteFromSnapshotPayload{
		ProductID:         in.ProductID,
		Cycle:             in.Cycle,
		KeypairID:         in.KeypairID,
		Name:              in.Name,
		Description:       strPtr(in.Description),
		SSHAndConsoleUser: in.SSHAndConsoleUser,
		ConsolePassword:   in.ConsolePassword,
		Promocode:         strPtr(in.Promocode),
		PayInvoiceWithCC:  yesNo(in.PayWithCreditCard),
	})
	if err != nil {
		return infer.CreateResponse[NeoliteVmFromSnapshotState]{}, err
	}
	if billing.AccountID == "" {
		return infer.CreateResponse[NeoliteVmFromSnapshotState]{},
			fmt.Errorf("create neolite vm from snapshot response tidak ada account_id: order_id=%s", billing.OrderID)
	}
	id := billing.AccountID
	partial := NeoliteVmFromSnapshotState{NeoliteVmFromSnapshotArgs: in, OrderID: billing.OrderID}

	acc, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (*client.AccountResource, error) {
			n, err := parseNeoID(id)
			if err != nil {
				return nil, err
			}
			return c.Neolite().AccountGet(ctx, n)
		},
		neoliteStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return infer.CreateResponse[NeoliteVmFromSnapshotState]{ID: id, Output: partial},
				infer.ResourceInitFailedError{Reasons: []string{
					fmt.Sprintf("neolite vm from snapshot %s not active yet: %s", id, err),
				}}
		}
		return infer.CreateResponse[NeoliteVmFromSnapshotState]{}, err
	}

	state := partial
	state.Status = acc.Status
	return infer.CreateResponse[NeoliteVmFromSnapshotState]{ID: id, Output: state}, nil
}

func (NeoliteVmFromSnapshot) Read(
	ctx context.Context, req infer.ReadRequest[NeoliteVmFromSnapshotArgs, NeoliteVmFromSnapshotState],
) (infer.ReadResponse[NeoliteVmFromSnapshotArgs, NeoliteVmFromSnapshotState], error) {
	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.ReadResponse[NeoliteVmFromSnapshotArgs, NeoliteVmFromSnapshotState]{}, err
	}
	acc, err := c.Neolite().AccountGet(ctx, id)
	if err != nil {
		if client.IsNotFound(err) {
			return infer.ReadResponse[NeoliteVmFromSnapshotArgs, NeoliteVmFromSnapshotState]{},
				fmt.Errorf("neolite vm from snapshot %s not found", req.ID)
		}
		return infer.ReadResponse[NeoliteVmFromSnapshotArgs, NeoliteVmFromSnapshotState]{}, err
	}
	state := req.State
	state.Status = acc.Status
	return infer.ReadResponse[NeoliteVmFromSnapshotArgs, NeoliteVmFromSnapshotState]{State: state}, nil
}

func (NeoliteVmFromSnapshot) Delete(
	ctx context.Context, req infer.DeleteRequest[NeoliteVmFromSnapshotState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	// delete = delete VM hasil restore.
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if _, err := c.Neolite().VMDelete(ctx, id); err != nil {
		if client.IsNotFound(err) {
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
