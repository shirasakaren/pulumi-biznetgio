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
