package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type NeoliteProService struct{ client *Client }

func (s *NeoliteProService) path(seg, status string) string {
	return withStatus("/neolite-pros"+seg, status)
}

func (s *NeoliteProService) AccountsList(ctx context.Context, status string) ([]AccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]AccountResource](raw)
}

func (s *NeoliteProService) AccountGet(ctx context.Context, accountID int64) (*AccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*AccountResource](raw)
}

func (s *NeoliteProService) VMCreate(ctx context.Context, req NeoliteCreatePayload) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolite-pros", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteProService) VMDetails(ctx context.Context, accountID int64) (*VirtualMachineResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/accounts/%d/vm-details", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*VirtualMachineResource](raw)
}

func (s *NeoliteProService) VMDelete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolite-pros/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) VMState(ctx context.Context, accountID int64, state string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolite-pros/accounts/%d/vm-state/%s", accountID, url.PathEscape(state)), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) VMRebuild(ctx context.Context, accountID int64, osName string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolite-pros/accounts/%d/rebuild", accountID),
		map[string]string{"select_os": osName})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) VMChangeName(ctx context.Context, accountID int64, name string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolite-pros/accounts/%d/change-vm-name", accountID),
		map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) VMChangeKeypair(ctx context.Context, accountID, keypairID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolite-pros/accounts/%d/keypair", accountID),
		map[string]int64{"keypair_id": keypairID})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) VMChangePackage(
	ctx context.Context, accountID int64, req ChangePackagePayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neolite-pros/accounts/%d/change-package", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteProService) VMChangeStorage(
	ctx context.Context, accountID int64, req UpgradePayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolite-pros/accounts/%d/storage", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteProService) ChangePackageOptions(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neolite-pros/accounts/%d/change-package", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) StorageOptions(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neolite-pros/accounts/%d/storage", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) DiskCreate(ctx context.Context, req NeoliteDiskCreatePayload) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolite-pros/disks", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) DisksList(ctx context.Context, status string) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/disks/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *NeoliteProService) DiskGet(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/disks/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) DiskUpgrade(
	ctx context.Context, accountID int64, req DiskUpgradePayload,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut, fmt.Sprintf("/neolite-pros/disks/accounts/%d", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) DiskDelete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolite-pros/disks/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) DiskProductList(ctx context.Context) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolite-pros/disks/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *NeoliteProService) DiskProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/disks/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) AccountSnapshotList(ctx context.Context, status string) ([]SnapshotAccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/snapshots/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]SnapshotAccountResource](raw)
}

func (s *NeoliteProService) AccountSnapshotGet(ctx context.Context, accountID int64) (*SnapshotAccountResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/snapshots/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*SnapshotAccountResource](raw)
}

func (s *NeoliteProService) SnapshotCreate(
	ctx context.Context, accountID int64, req SnapshotPayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neolite-pros/accounts/%d/snapshot", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteProService) SnapshotRestore(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neolite-pros/snapshots/accounts/%d/restore", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) SnapshotRestoreWith(
	ctx context.Context, accountID int64, req NeoliteFromSnapshotPayload,
) (*BillingResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neolite-pros/snapshots/accounts/%d/create", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*BillingResource](raw)
}

func (s *NeoliteProService) SnapshotDelete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolite-pros/snapshots/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) SnapshotProductsList(ctx context.Context) ([]PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolite-pros/snapshots/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]PlanResource](raw)
}

func (s *NeoliteProService) SnapshotProductGet(ctx context.Context, productID int64) (*PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/snapshots/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*PlanResource](raw)
}

func (s *NeoliteProService) KeypairList(ctx context.Context) ([]KeypairResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolite-pros/keypairs/", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]KeypairResource](raw)
}

func (s *NeoliteProService) KeypairCreate(ctx context.Context, name string) (*KeypairResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolite-pros/keypairs/", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	return decodeJSON[*KeypairResource](raw)
}

// KeypairCreateRaw creates a keypair and returns the raw response map - private_key only
// appears in this response and the field is undocumented, so don't decode it typed.
func (s *NeoliteProService) KeypairCreateRaw(ctx context.Context, name string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolite-pros/keypairs/", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) KeypairImport(ctx context.Context, req KeypairImportPayload) (*KeypairResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neolite-pros/keypairs/import", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*KeypairResource](raw)
}

func (s *NeoliteProService) KeypairDelete(ctx context.Context, keypairID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neolite-pros/keypairs/%d", keypairID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *NeoliteProService) ProductList(ctx context.Context) ([]PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neolite-pros/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]PlanResource](raw)
}

func (s *NeoliteProService) ProductGet(ctx context.Context, productID int64) (*PlanResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*PlanResource](raw)
}

func (s *NeoliteProService) ProductOSList(ctx context.Context, productID int64) ([]OsResource, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neolite-pros/products/%d/oss", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]OsResource](raw)
}

func (s *NeoliteProService) ProductIPAvailability(ctx context.Context, productID int64) (*IPAvailability, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neolite-pros/products/%d/ip-availability", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[*IPAvailability](raw)
}
