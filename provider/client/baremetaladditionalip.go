package client

import (
	"context"
	"fmt"
	"net/http"
)

type BaremetalAdditionalIPService struct{ client *Client }

func (s *BaremetalAdditionalIPService) path(seg, status string) string {
	return withStatus("/baremetal-additional-ips"+seg, status)
}

func (s *BaremetalAdditionalIPService) Create(
	ctx context.Context, req AdditionalIPCreatePayload,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/baremetal-additional-ips", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) List(ctx context.Context, status string) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

// docs bilang id string di get, int di delete — ikut docs aja
func (s *BaremetalAdditionalIPService) Get(ctx context.Context, accountID string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/baremetal-additional-ips/%s", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) Delete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/baremetal-additional-ips/%d", accountID), nil)
	if err != nil {
		return nil, err
