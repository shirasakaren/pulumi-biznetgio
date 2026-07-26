package client

import (
	"context"
	"fmt"
	"net/http"
)

type BaremetalElasticStorageService struct{ client *Client }

func (s *BaremetalElasticStorageService) path(seg, status string) string {
	return withStatus("/baremetal-neo-elastic-storages"+seg, status)
}

func (s *BaremetalElasticStorageService) Create(
	ctx context.Context, req NeoElasticStorageCreatePayload,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/baremetal-neo-elastic-storages", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalElasticStorageService) List(ctx context.Context, status string) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *BaremetalElasticStorageService) Get(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/baremetal-neo-elastic-storages/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalElasticStorageService) Upgrade(
	ctx context.Context, accountID int64, req UpgradeNeoElasticStorage,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/baremetal-neo-elastic-storages/%d", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalElasticStorageService) ChangePackage(
	ctx context.Context, accountID int64, req ChangePackageNeoElasticStorage,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/baremetal-neo-elastic-storages/%d", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalElasticStorageService) Delete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete,
		fmt.Sprintf("/baremetal-neo-elastic-storages/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

// wip 768
