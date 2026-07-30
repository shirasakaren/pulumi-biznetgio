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
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) AssignmentsList(ctx context.Context, accountID int64) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/baremetal-additional-ips/%d/assigns", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) AssignToMachine(
	ctx context.Context, accountID, metalAccountID int64,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/baremetal-additional-ips/%d/assigns", accountID),
		map[string]int64{"metal_account_id": metalAccountID})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) AssignmentGet(
	ctx context.Context, accountID, metalAccountID int64,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/baremetal-additional-ips/%d/assigns/%d", accountID, metalAccountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) Unassign(
	ctx context.Context, accountID, metalAccountID int64,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete,
		fmt.Sprintf("/baremetal-additional-ips/%d/assigns/%d", accountID, metalAccountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) AssignmentsByMetal(
	ctx context.Context, metalAccountID int64,
) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/baremetal-additional-ips/assignments-by-metal-account-id/%d", metalAccountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) ProductList(ctx context.Context) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/baremetal-additional-ips/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) ProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/baremetal-additional-ips/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalAdditionalIPService) Regions(ctx context.Context) ([]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/baremetal-additional-ips/regions", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]any](raw)
}

type AdditionalIPCreatePayload struct {
	ProductID        int64  `json:"product_id"`
	Cycle            string `json:"cycle"`
	Region           string `json:"region,omitempty"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}
