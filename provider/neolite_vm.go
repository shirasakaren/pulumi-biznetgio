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

type NeoliteVm struct{}

type NeoliteVmArgs struct {
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
	MigrateToPro      *string `pulumi:"migrateToPro,optional"`
	DiskSize          *int64  `pulumi:"diskSize,optional"`
}

type NeoliteLastInvoice struct {
	ID          int64  `pulumi:"id"`
	PaidID      int64  `pulumi:"paidId"`
	Status      string `pulumi:"status"`
	Date        string `pulumi:"date"`
	Duedate     string `pulumi:"duedate"`
	Paybefore   string `pulumi:"paybefore"`
	Datepaid    string `pulumi:"datepaid"`
