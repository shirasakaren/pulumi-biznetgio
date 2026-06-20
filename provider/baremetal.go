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
