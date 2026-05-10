package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/biznetgio/pulumi-biznetgio/provider/client"
)

type BaremetalElasticStorage struct{}

type BaremetalElasticStorageArgs struct {
	ProductID         int64   `pulumi:"productId"`
	Cycle             string  `pulumi:"cycle"`
	StorageName       string  `pulumi:"storageName"`
	MetalAccountID    int64   `pulumi:"metalAccountId"`
	Size              *int64  `pulumi:"size,optional"`
	Promocode         *string `pulumi:"promocode,optional"`
	PayWithCreditCard *bool   `pulumi:"payWithCreditCard,optional"`
}

type BaremetalElasticStorageState struct {
	BaremetalElasticStorageArgs
	Status    string  `pulumi:"status"`
	CreatedAt *string `pulumi:"createdAt"`
	Raw       string  `pulumi:"raw" provider:"secret"`
}

func (a *BaremetalElasticStorageArgs) Annotate(ann infer.Annotator) {
	ann.Describe(&a.ProductID, "Product id from `GET /baremetal-neo-elastic-storages/products`. "+
		"Changing it triggers change-package (`POST .../{account_id}`).")
	ann.Describe(&a.Cycle, "Billing cycle, e.g. `m` for monthly or `a` for annual. "+
		"Create-only, changing it replaces the storage.")
	ann.Describe(&a.StorageName, "Name of the storage. Create-only, changing it replaces the storage.")
	ann.Describe(&a.MetalAccountID, "Account id of the target baremetal. "+
		"Can only be set at creation (no re-attach endpoint).")
	ann.Describe(&a.Size, "Storage size in GB. Defaults to 100. Changing it triggers upgrade "+
		"(`PUT .../{account_id}`) — grow-only, shrinking is rejected by the API.")
	ann.SetDefault(&a.Size, int64(100))
	ann.Describe(&a.Promocode, "Promo code to apply at creation.")
