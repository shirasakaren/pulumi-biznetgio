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

func (s *ObjectStorageService) BucketCreate(
	ctx context.Context, accountID int64, req NOSCreateBucketPayload,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		fmt.Sprintf("/object-storages/accounts/%d/buckets", accountID), req)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) BucketSetACL(
	ctx context.Context, accountID int64, bucketName, acl string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		s.bucketPath(accountID, bucketName, ""), map[string]string{"acl": acl})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) BucketDelete(
	ctx context.Context, accountID int64, bucketName string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete, s.bucketPath(accountID, bucketName, ""), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) BucketUsage(
	ctx context.Context, accountID int64, bucketName string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.bucketPath(accountID, bucketName, "/usage"), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) Info(
	ctx context.Context, accountID int64, bucketName, objectOrDirectory string,
) (map[string]any, error) {
	p := s.bucketPath(accountID, bucketName, "/info")
	if objectOrDirectory != "" {
		p += "?object_or_directory=" + url.QueryEscape(objectOrDirectory)
	}
	raw, err := s.client.doJSON(ctx, http.MethodGet, p, nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) ObjectsList(
	ctx context.Context, accountID int64, bucketName string,
) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet, s.bucketPath(accountID, bucketName, "/objects"), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *ObjectStorageService) ObjectsListInDirectory(
	ctx context.Context, accountID int64, bucketName, directory string,
) ([]map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodGet,
		s.bucketPath(accountID, bucketName, "/objects/"+url.PathEscape(directory)+"/"), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[[]map[string]any](raw)
}

func (s *ObjectStorageService) DirectoryCreate(
	ctx context.Context, accountID int64, bucketName, directory string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		s.bucketPath(accountID, bucketName, "/objects/"+url.PathEscape(directory)+"/"), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

// responnya bytes, bukan envelope — jangan di-unwrap
func (s *ObjectStorageService) ObjectDownload(
	ctx context.Context, accountID int64, bucketName, objectName string,
) ([]byte, error) {
	return s.client.do(ctx, http.MethodGet,
		s.bucketPath(accountID, bucketName, "/objects/"+url.PathEscape(objectName)),
		"application/json", nil)
}

func (s *ObjectStorageService) ObjectSetACL(
	ctx context.Context, accountID int64, bucketName, objectOrDirectory, acl string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPut,
		s.bucketPath(accountID, bucketName, "/objects/"+url.PathEscape(objectOrDirectory)),
		map[string]string{"acl": acl})
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) ObjectDelete(
	ctx context.Context, accountID int64, bucketName, objectOrDirectory string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodDelete,
		s.bucketPath(accountID, bucketName, "/objects/"+url.PathEscape(objectOrDirectory)), nil)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) Upload(
	ctx context.Context, accountID int64, bucketName, directory, filename string, content []byte,
) (map[string]any, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := fw.Write(content); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	p := s.bucketPath(accountID, bucketName, "/upload")
	if directory != "" {
		p += "?directory=" + url.QueryEscape(directory)
	}
	raw, err := s.client.do(ctx, http.MethodPost, p, w.FormDataContentType(), buf.Bytes())
	if err != nil {
		return nil, err
	}
	payload, err := unwrapJSON(raw)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](payload)
}

func (s *ObjectStorageService) ObjectURL(
	ctx context.Context, accountID int64, bucketName, objectName string, expiry int64,
) (map[string]any, error) {
	var body any
	if expiry != 0 {
		body = map[string]int64{"expiry": expiry}
	}
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		s.bucketPath(accountID, bucketName, "/url/"+url.PathEscape(objectName)), body)
	if err != nil {
		return nil, err
	}
	return decodeJSON[map[string]any](raw)
}

func (s *ObjectStorageService) ObjectCopy(
	ctx context.Context, accountID int64, bucketName, toBucketName, objectName string,
) (map[string]any, error) {
	raw, err := s.client.doJSON(ctx, http.MethodPost,
		s.bucketPath(accountID, bucketName, "/copy/"+toBucketName+"/"+objectName), nil)
	if err != nil {
