package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/shirasakaren/pulumi-biznetgio/provider/client"
)

type NeoliteProSnapshot struct{}

type NeoliteProSnapshotArgs struct {
	NeoliteAccountID  int64   `pulumi:"neoliteAccountId"`
	Name              *string `pulumi:"name,optional"`
	Description       *string `pulumi:"description,optional"`
	Cycle             string  `pulumi:"cycle"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
}

type NeoliteProSnapshotState struct {
	NeoliteProSnapshotArgs
	OrderID string `pulumi:"orderId"`
	Status  string `pulumi:"status"`
}

func (a *NeoliteProSnapshotArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.NeoliteAccountID, "Account id of the NEO Lite Pro VM being snapshotted.")
	ann.Describe(&a.Name, "Snapshot name. Defaults to `snapshot-name`.")
	ann.SetDefault(&a.Name, "snapshot-name")
	ann.Describe(&a.Description, "Snapshot description.")
	ann.SetDefault(&a.Description, "")
	ann.Describe(&a.Cycle, "Billing cycle: m monthly, a annual, q quarterly, s semiannual, "+
		"b biennial, t triennial, p4, p5.")
	ann.Describe(&a.PayWithCreditCard, "Pay the invoice with the registered credit card at order time. "+
		"Defaults to true (auto-charge); set false to leave it unpaid in the portal until settled.")
	ann.SetDefault(&a.PayWithCreditCard, true)
	ann.Describe(&a.Promocode, "Promo code to apply at order.")
	ann.SetDefault(&a.Promocode, "")
}

func (s *NeoliteProSnapshotState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.OrderID, "Order id from the creation response.")
	ann.Describe(&s.Status, "Snapshot status (Active, Pending, Suspended, Terminated).")
}

func (NeoliteProSnapshot) Create(
	ctx context.Context, req infer.CreateRequest[NeoliteProSnapshotArgs],
) (infer.CreateResponse[NeoliteProSnapshotState], error) {
	if req.DryRun {
		return infer.CreateResponse[NeoliteProSnapshotState]{
			ID:     "0",
			Output: NeoliteProSnapshotState{NeoliteProSnapshotArgs: req.Inputs},
		}, nil
	}

	c := GetClient(ctx)
	in := req.Inputs
	billing, err := c.NeolitePro().SnapshotCreate(ctx, in.NeoliteAccountID, client.SnapshotPayload{
		Cycle:            in.Cycle,
		Name:             strPtr(in.Name),
		Description:      strPtr(in.Description),
		Promocode:        strPtr(in.Promocode),
		PayInvoiceWithCC: yesNo(in.PayWithCreditCard),
	})
	if err != nil {
		return infer.CreateResponse[NeoliteProSnapshotState]{}, err
	}
	if billing.AccountID == "" {
		return infer.CreateResponse[NeoliteProSnapshotState]{},
			fmt.Errorf("create neolite pro snapshot response tidak ada account_id: order_id=%s", billing.OrderID)
	}
	id := billing.AccountID
	partial := NeoliteProSnapshotState{NeoliteProSnapshotArgs: in, OrderID: billing.OrderID}

	acc, err := client.WaitForStatus(ctx, 5*time.Second,
		func(ctx context.Context) (*client.SnapshotAccountResource, error) {
			n, err := parseNeoID(id)
			if err != nil {
				return nil, err
			}
			return c.NeolitePro().AccountSnapshotGet(ctx, n)
		},
		neoliteSnapshotStatus, []string{"active"}, []string{"suspended", "terminated"})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return infer.CreateResponse[NeoliteProSnapshotState]{ID: id, Output: partial},
				infer.ResourceInitFailedError{Reasons: []string{fmt.Sprintf("neolite pro snapshot %s not active yet: %s", id, err)}}
		}
		return infer.CreateResponse[NeoliteProSnapshotState]{}, err
	}

	state := partial
	state.Status = acc.Status
	return infer.CreateResponse[NeoliteProSnapshotState]{ID: id, Output: state}, nil
}

func (NeoliteProSnapshot) Read(
	ctx context.Context, req infer.ReadRequest[NeoliteProSnapshotArgs, NeoliteProSnapshotState],
) (infer.ReadResponse[NeoliteProSnapshotArgs, NeoliteProSnapshotState], error) {
	c := GetClient(ctx)
	list, err := c.NeolitePro().AccountSnapshotList(ctx, "")
	if err != nil {
		return infer.ReadResponse[NeoliteProSnapshotArgs, NeoliteProSnapshotState]{}, err
	}
	for _, sn := range list {
		if sn.AccountID == req.ID {
			state := req.State
			state.Status = sn.Status
			if sn.ExtraDetails.Name != "" {
				state.Name = &sn.ExtraDetails.Name
			}
			if sn.ExtraDetails.Description != "" {
				state.Description = &sn.ExtraDetails.Description
			}
			return infer.ReadResponse[NeoliteProSnapshotArgs, NeoliteProSnapshotState]{State: state}, nil
		}
	}
	return infer.ReadResponse[NeoliteProSnapshotArgs, NeoliteProSnapshotState]{},
		fmt.Errorf("neolite pro snapshot %s not found", req.ID)
}

func (NeoliteProSnapshot) Delete(
	ctx context.Context, req infer.DeleteRequest[NeoliteProSnapshotState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	id, err := parseNeoID(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if _, err := c.NeolitePro().SnapshotDelete(ctx, id); err != nil {
		if client.IsNotFound(err) {
			return infer.DeleteResponse{}, nil
		}
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}
