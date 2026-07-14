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
