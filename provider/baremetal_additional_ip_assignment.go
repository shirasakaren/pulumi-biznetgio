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
