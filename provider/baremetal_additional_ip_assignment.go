package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type BaremetalAdditionalIpAssignment struct{}

type BaremetalAdditionalIpAssignmentArgs struct {
	AdditionalIPID int64 `pulumi:"additionalIpId"`
	MetalAccountID int64 `pulumi:"metalAccountId"`
}

type BaremetalAdditionalIpAssignmentState struct {
	BaremetalAdditionalIpAssignmentArgs
	Status string `pulumi:"status"`
	Raw    string `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalAdditionalIpAssignmentArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.AdditionalIPID, "Account id of the additional IP to attach, from BaremetalAdditionalIp.")
	ann.Describe(&a.MetalAccountID, "Account id of the baremetal to attach to. "+
		"Changing the target replaces the assignment.")
}

func (s *BaremetalAdditionalIpAssignmentState) Annotate(ann infer.Annotator) {
	ann.Describe(&s.Status, "Assignment status, when present in the response.")
	ann.Describe(&s.Raw, "Raw JSON of the last read response, for anything not modeled yet.")
}

func (BaremetalAdditionalIpAssignment) Create(
	ctx context.Context, req infer.CreateRequest[BaremetalAdditionalIpAssignmentArgs],
) (infer.CreateResponse[BaremetalAdditionalIpAssignmentState], error) {
	resp := infer.CreateResponse[BaremetalAdditionalIpAssignmentState]{
		Output: assignmentStateFromMap(ctx, req.Inputs, nil),
	}
	if req.DryRun {
		return resp, nil
	}
	c := GetClient(ctx)
	a := req.Inputs

	raw, err := c.BaremetalAdditionalIP().AssignToMachine(ctx, a.AdditionalIPID, a.MetalAccountID)
	if err != nil {
		return infer.CreateResponse[BaremetalAdditionalIpAssignmentState]{}, err
	}
	resp.ID = fmt.Sprintf("%d:%d", a.MetalAccountID, a.AdditionalIPID)
	resp.Output = assignmentStateFromMap(ctx, a, raw)
	return resp, nil
}

func (BaremetalAdditionalIpAssignment) Read(
	ctx context.Context, req infer.ReadRequest[BaremetalAdditionalIpAssignmentArgs, BaremetalAdditionalIpAssignmentState],
) (infer.ReadResponse[BaremetalAdditionalIpAssignmentArgs, BaremetalAdditionalIpAssignmentState], error) {
	resp := infer.ReadResponse[BaremetalAdditionalIpAssignmentArgs, BaremetalAdditionalIpAssignmentState]{
		ID:     req.ID,
		Inputs: req.Inputs,
		State:  assignmentStateFromMap(ctx, req.Inputs, nil),
	}
	c := GetClient(ctx)
	ipID, metalID, err := parseAssignmentID(req.ID)
	if err != nil {
		return resp, err
	}
	m, err := c.BaremetalAdditionalIP().AssignmentGet(ctx, ipID, metalID)
	if err != nil {
		if client.IsNotFound(err) {
			return resp, fmt.Errorf("biznetgio: assignment %s not found", req.ID)
		}
		return resp, err
	}
	resp.State = assignmentStateFromMap(ctx, req.Inputs, m)
	return resp, nil
}

func (BaremetalAdditionalIpAssignment) Delete(
	ctx context.Context, req infer.DeleteRequest[BaremetalAdditionalIpAssignmentState],
) (infer.DeleteResponse, error) {
	c := GetClient(ctx)
	ipID, metalID, err := parseAssignmentID(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if _, err := c.BaremetalAdditionalIP().Unassign(ctx, ipID, metalID); err != nil && !client.IsNotFound(err) {
		return infer.DeleteResponse{}, err
	}
	return infer.DeleteResponse{}, nil
}

func assignmentStateFromMap(
	_ context.Context, args BaremetalAdditionalIpAssignmentArgs, m map[string]any,
) BaremetalAdditionalIpAssignmentState {
	st := BaremetalAdditionalIpAssignmentState{BaremetalAdditionalIpAssignmentArgs: args}
	if m == nil {
		return st
	}
	st.Status = bmtStringDefault(m, "status", "state")
	st.Raw = string(bmtJSON(m))
	return st
}

// parseAssignmentID splits the composite id `<metal_account_id>:<additional_ip_id>`.
func parseAssignmentID(id string) (ipID, metalID int64, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("biznetgio: invalid assignment id %q, want `<metal_account_id>:<additional_ip_id>`", id)
	}
	metal, err1 := strconv.ParseInt(parts[0], 10, 64)
	ip, err2 := strconv.ParseInt(parts[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("biznetgio: invalid assignment id %q", id)
	}
	return ip, metal, nil
}
