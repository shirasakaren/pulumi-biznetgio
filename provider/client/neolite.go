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
