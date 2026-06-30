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
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) Delete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/neo-gpus/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) VMStatusGet(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neo-gpus/accounts/%d/vm-status", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) VMStatusSet(ctx context.Context, accountID int64, status string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neo-gpus/accounts/%d/vm-status", accountID),
		map[string]string{"status": status})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) Rebuild(ctx context.Context, accountID int64, selectOS string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/neo-gpus/accounts/%d/rebuild", accountID),
		map[string]string{"select_os": selectOS})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) ReserveAdditionalHours(ctx context.Context, accountID, hours int64) (map[string]any, error) {
	var body any
	if hours != 0 {
		body = map[string]int64{"hours": hours}
	}
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neo-gpus/accounts/%d/reserve-additional-hours", accountID), body)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) ConsoleAccess(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/neo-gpus/accounts/%d/console-access", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) GraphMonitor(ctx context.Context, accountID int64, timeframe string) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/neo-gpus/accounts/%d/graph-monitor?timeframe=%s", accountID, url.QueryEscape(timeframe)), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *GPUService) KeypairList(ctx context.Context) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/neo-gpus/keypairs/", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

