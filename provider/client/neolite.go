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
