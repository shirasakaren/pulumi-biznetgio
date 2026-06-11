package client

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

type ObjectStorageService struct{ client *Client }

func (s *ObjectStorageService) path(seg, status string) string {
	return withStatus("/object-storages"+seg, status)
}

func (s *ObjectStorageService) bucketPath(accountID int64, bucketName, rest string) string {
	return fmt.Sprintf("/object-storages/accounts/%d/buckets/%s%s", accountID, bucketName, rest)
}

func (s *ObjectStorageService) Create(ctx context.Context, req NOSCreatePayload) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost, "/object-storages", req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) AccountsList(ctx context.Context, status string) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.path("/accounts", status), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *ObjectStorageService) AccountGet(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/object-storages/accounts/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) QuotaUpgrade(
	ctx context.Context, accountID int64, req NOSQuotaUpgradePayload,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/object-storages/accounts/%d", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) Delete(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/object-storages/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) ProductList(ctx context.Context) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, "/object-storages/products", nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *ObjectStorageService) ProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, fmt.Sprintf("/object-storages/products/%d", productID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) CredentialsList(ctx context.Context, accountID int64) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/object-storages/accounts/%d/credentials", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *ObjectStorageService) CredentialCreate(ctx context.Context, accountID int64) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/object-storages/accounts/%d/credentials", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) CredentialUpdateStatus(
	ctx context.Context, accountID int64, accessKey string, active bool,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		fmt.Sprintf("/object-storages/accounts/%d/credentials/%s", accountID, accessKey),
		map[string]bool{"active": active})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) CredentialDelete(
	ctx context.Context, accountID int64, accessKey string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete,
		fmt.Sprintf("/object-storages/accounts/%d/credentials/%s", accountID, accessKey), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) BucketsList(ctx context.Context, accountID int64) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		fmt.Sprintf("/object-storages/accounts/%d/buckets", accountID), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}
