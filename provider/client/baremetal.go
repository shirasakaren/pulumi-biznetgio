package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type BaremetalService struct{ client *Client }

func (s *BaremetalService) path(seg, status string) string {
	return withStatus("/baremetals"+seg, status)
}

func (s *BaremetalService) Create(ctx context.Context, req BaremetalCreatePayload) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/baremetals", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) AccountsList(ctx context.Context, status string) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *BaremetalService) AccountGet(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/baremetals/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) UpdateLabel(ctx context.Context, accountID int64, label string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut, fmt.Sprintf("/baremetals/%d", accountID),
		map[string]string{"label": label})
	if err != nil {
		return nil, err
	}
