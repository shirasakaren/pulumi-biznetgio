package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type NeoliteService struct{ client *Client }

func (s *NeoliteService) path(seg, status string) string {
	return withStatus("/neolites"+seg, status)
}

func (s *NeoliteService) AccountsList(ctx context.Context, status string) ([]AccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]AccountResource](raw)
}

func (s *NeoliteService) AccountGet(ctx context.Context, accountID int64) (*AccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*AccountResource](raw)
}

func (s *NeoliteService) VMCreate(ctx context.Context, req NeoliteCreatePayload) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolites", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteService) VMDetails(ctx context.Context, accountID int64) (*VirtualMachineResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/accounts/%d/vm-details", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*VirtualMachineResource](raw)
}

func (s *NeoliteService) VMDelete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolites/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) VMState(ctx context.Context, accountID int64, state string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolites/accounts/%d/vm-state/%s", accountID, url.PathEscape(state)), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) VMRebuild(ctx context.Context, accountID int64, osName string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolites/accounts/%d/rebuild", accountID),
		map[string]string{"select_os": osName})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) VMChangeName(ctx context.Context, accountID int64, name string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolites/accounts/%d/change-vm-name", accountID),
		map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) VMChangeKeypair(ctx context.Context, accountID, keypairID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolites/accounts/%d/keypair", accountID),
		map[string]int64{"keypair_id": keypairID})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) VMChangePackage(
	ctx context.Context, accountID int64, req ChangePackagePayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neolites/accounts/%d/change-package", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteService) VMChangeStorage(
	ctx context.Context, accountID int64, req UpgradePayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolites/accounts/%d/storage", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteService) ChangePackageOptions(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neolites/accounts/%d/change-package", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) StorageOptions(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neolites/accounts/%d/storage", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) MigrateToPro(
	ctx context.Context, accountID int64, req MigrateToProPayload,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neolites/accounts/%d/migrate-to-pro", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) MigrateToProProducts(ctx context.Context, accountID int64) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neolites/accounts/%d/migrate-to-pro/products", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *NeoliteService) DiskCreate(ctx context.Context, req NeoliteDiskCreatePayload) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolites/disks", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) DisksList(ctx context.Context, status string) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/disks/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *NeoliteService) DiskGet(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/disks/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) DiskUpgrade(
	ctx context.Context, accountID int64, req DiskUpgradePayload,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut, fmt.Sprintf("/neolites/disks/accounts/%d", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) DiskDelete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolites/disks/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) DiskProductList(ctx context.Context) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolites/disks/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *NeoliteService) DiskProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/disks/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) AccountSnapshotList(ctx context.Context, status string) ([]SnapshotAccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/snapshots/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]SnapshotAccountResource](raw)
}

func (s *NeoliteService) AccountSnapshotGet(ctx context.Context, accountID int64) (*SnapshotAccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/snapshots/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*SnapshotAccountResource](raw)
}

func (s *NeoliteService) SnapshotCreate(
	ctx context.Context, accountID int64, req SnapshotPayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neolites/accounts/%d/snapshot", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteService) SnapshotRestore(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolites/snapshots/accounts/%d/restore", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) SnapshotRestoreWith(
	ctx context.Context, accountID int64, req NeoliteFromSnapshotPayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neolites/snapshots/accounts/%d/create", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteService) SnapshotDelete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolites/snapshots/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) SnapshotProductsList(ctx context.Context) ([]PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolites/snapshots/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]PlanResource](raw)
}

func (s *NeoliteService) SnapshotProductGet(ctx context.Context, productID int64) (*PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/snapshots/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*PlanResource](raw)
}

func (s *NeoliteService) KeypairList(ctx context.Context) ([]KeypairResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolites/keypairs/", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]KeypairResource](raw)
}

func (s *NeoliteService) KeypairCreate(ctx context.Context, name string) (*KeypairResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolites/keypairs/", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	return decodeJSON[*KeypairResource](raw)
}

// KeypairCreateRaw create keypair, return raw response map — private key cuma
// keluar di response ini dan field-nya undocumented, jadi jangan di-decode typed.
func (s *NeoliteService) KeypairCreateRaw(ctx context.Context, name string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolites/keypairs/", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) KeypairImport(ctx context.Context, req KeypairImportPayload) (*KeypairResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolites/keypairs/import", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*KeypairResource](raw)
}

func (s *NeoliteService) KeypairDelete(ctx context.Context, keypairID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolites/keypairs/%d", keypairID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteService) ProductList(ctx context.Context) ([]PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolites/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]PlanResource](raw)
}

func (s *NeoliteService) ProductGet(ctx context.Context, productID int64) (*PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*PlanResource](raw)
}

func (s *NeoliteService) ProductOSList(ctx context.Context, productID int64) ([]OsResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolites/products/%d/oss", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]OsResource](raw)
}

func (s *NeoliteService) ProductIPAvailability(ctx context.Context, productID int64) (*IPAvailability, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neolites/products/%d/ip-availability", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*IPAvailability](raw)
}

type BillingResource struct {
	OrderID   string `json:"order_id"`
	AccountID string `json:"account_id"`
}

type AccountResource struct {
	AccountID       string       `json:"account_id"`
	Domain          string       `json:"domain"`
	Status          string       `json:"status"`
	Billingcycle    string       `json:"billingcycle"`
	DateCreated     string       `json:"date_created"`
	NextDue         string       `json:"next_due"`
	RecurringAmount int64        `json:"recurring_amount"`
	ExtraDetails    ExtraDetails `json:"extra_details"`
	ProductID       int64        `json:"product_id"`
	ProductName     string       `json:"product_name"`
	Description     string       `json:"description"`
	CategoryID      int64        `json:"category_id"`
	CategoryName    string       `json:"category_name"`
	LastInvoice     LastInvoice  `json:"last_invoice"`
}

type LastInvoice struct {
	ID          int64  `json:"id"`
	PaidID      int64  `json:"paid_id"`
	Status      string `json:"status"`
	Date        string `json:"date"`
	Duedate     string `json:"duedate"`
	Paybefore   string `json:"paybefore"`
	Datepaid    string `json:"datepaid"`
	InvoiceType string `json:"invoice_type"`
}

type ExtraDetails struct {
	Region      string  `json:"region"`
	RegionLabel string  `json:"region_label"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
	TenantID    *string `json:"tenant_id"`
	CIUser      string  `json:"ciuser"`
	CIPassword  string  `json:"cipassword"`
	KeypairID   int64   `json:"neosshkey_id"`
	SSHKeys     string  `json:"sshkeys"`
	OSName      string  `json:"osname"`
	DiskSize    string  `json:"disk_size"`
}

type VirtualMachineResource struct {
	VMID    int64  `json:"vmid"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Uptime  int64  `json:"uptime"`
	MaxDisk int64  `json:"maxdisk"`
	MaxMem  int64  `json:"maxmem"`
	Mem     int64  `json:"mem"`
	CPUs    int64  `json:"cpus"`
}

type KeypairResource struct {
	KeypairID int64  `json:"keypair_id"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

// UnmarshalJSON toleran: keypair_id bisa number, string numerik, atau malah
// ke-rename jadi neosshkey_id (wire field beda-beda per endpoint).
func (k *KeypairResource) UnmarshalJSON(b []byte) error {
	var a struct {
		KeypairID   jsonInt64 `json:"keypair_id"`
		NeoSSHKeyID jsonInt64 `json:"neosshkey_id"`
		Name        string    `json:"name"`
		PublicKey   string    `json:"public_key"`
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	k.KeypairID = int64(a.KeypairID)
	if k.KeypairID == 0 {
		k.KeypairID = int64(a.NeoSSHKeyID)
	}
	k.Name = a.Name
	k.PublicKey = a.PublicKey
	return nil
}

type PlanResource struct {
	ProductID    int64     `json:"product_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CategoryID   int64     `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Options      Options   `json:"options"`
	Billing      []Billing `json:"billing"`
}

// UnmarshalJSON toleran: product_id juga bisa string numerik di wire.
func (p *PlanResource) UnmarshalJSON(b []byte) error {
	var a struct {
		ProductID    jsonInt64 `json:"product_id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		CategoryID   jsonInt64 `json:"category_id"`
		CategoryName string    `json:"category_name"`
		Options      Options   `json:"options"`
		Billing      []Billing `json:"billing"`
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	p.ProductID = int64(a.ProductID)
	p.Name = a.Name
	p.Description = a.Description
	p.CategoryID = int64(a.CategoryID)
	p.CategoryName = a.CategoryName
	p.Options = a.Options
	p.Billing = a.Billing
	return nil
}

type Options struct {
	Type           string `json:"type"`
	Cores          int64  `json:"cores"`
	Memory         int64  `json:"memory"`
	AllowDowngrade int64  `json:"allow_downgrade"`
}

type Billing struct {
	Label      string      `json:"label"`
	Cycle      string      `json:"cycle"`
	Price      int64       `json:"price"`
	Components []Component `json:"components"`
}

type Component struct {
	Label  string  `json:"label"`
	Field  string  `json:"field"`
	Prices []Price `json:"prices"`
}

type Price struct {
	QtyMin int64 `json:"qty_min"`
	QtyMax int64 `json:"qty_max"`
	Price  int64 `json:"price"`
}

type OsResource struct {
	VMID   int64  `json:"vmid"`
	Node   string `json:"node"`
	Name   string `json:"name"`
	MaxMem int64  `json:"maxmem"`
	MaxCPU int64  `json:"maxcpu"`
}

type IPAvailability struct {
	Available bool `json:"available"`
}

type SnapshotAccountResource struct {
	AccountID    string               `json:"account_id"`
	Status       string               `json:"status"`
	ExtraDetails SnapshotExtraDetails `json:"extra_details"`
}

type SnapshotExtraDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Region      string `json:"region"`
}

type SnapshotResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NeoliteCreatePayload struct {
	ProductID         int64  `json:"product_id"`
	Cycle             string `json:"cycle"`
	SelectOS          string `json:"select_os"`
	KeypairID         int64  `json:"keypair_id"`
	VMName            string `json:"vm_name,omitempty"`
	Description       string `json:"description,omitempty"`
