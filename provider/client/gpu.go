package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type GPUService struct{ client *Client }

func (s *GPUService) path(seg, status string) string {
	return withStatus("/neo-gpus"+seg, status)
}

func (s *GPUService) Create(ctx context.Context, req NEOGPUCreatePayload) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neo-gpus", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) CreateOneTime(ctx context.Context, req NEOGPUOneTimeCreatePayload) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/neo-gpus/one-time", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) AccountsList(ctx context.Context, status string) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *GPUService) AccountGet(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/neo-gpus/accounts/%d", accountID), nil)
// wip 432
