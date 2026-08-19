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
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) Delete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/baremetals/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) StateGet(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/baremetals/accounts/%d/state", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) StateSet(ctx context.Context, accountID int64, state string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/baremetals/accounts/%d/state/%s", accountID, url.PathEscape(state)), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) Rebuild(ctx context.Context, accountID int64, os string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut, fmt.Sprintf("/baremetals/%d/rebuild", accountID),
		map[string]string{"os": os})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) RebuildOSList(ctx context.Context, accountID int64) ([]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/baremetals/%d/rebuild/oss", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]any](raw)
}

func (s *BaremetalService) KeypairList(ctx context.Context) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/baremetals/keypairs/", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *BaremetalService) KeypairCreate(ctx context.Context, name string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/baremetals/keypairs/",
		map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) KeypairImport(ctx context.Context, req KeypairImportPayload) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/baremetals/keypairs/import", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) KeypairDelete(ctx context.Context, keypairID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/baremetals/keypairs/%d", keypairID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) OpenVPN(ctx context.Context) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/baremetals/openvpn", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) ProductList(ctx context.Context) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/baremetals/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *BaremetalService) ProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/baremetals/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *BaremetalService) ProductOSList(ctx context.Context, productID int64) ([]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/baremetals/products/%d/oss", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]any](raw)
}

func (s *BaremetalService) States(ctx context.Context) ([]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/baremetals/states", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]any](raw)
}

func (s *BaremetalService) AccountID(d map[string]any) string { return strField(d, "account_id") }
func (s *BaremetalService) ProductID(d map[string]any) int64  { return intField(d, "product_id") }

type BaremetalCreatePayload struct {
	ProductID        int64  `json:"product_id"`
	Cycle            string `json:"cycle"`
	SelectOS         string `json:"select_os,omitempty"`
	KeypairID        int64  `json:"keypair_id"`
	Label            string `json:"label"`
	PublicIP         int64  `json:"public_ip"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}
