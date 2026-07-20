package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.portal.biznetgio.com/v1"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string, timeout time.Duration) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) Neolite() *NeoliteService       { return &NeoliteService{client: c} }
func (c *Client) NeolitePro() *NeoliteProService { return &NeoliteProService{client: c} }
func (c *Client) Baremetal() *BaremetalService   { return &BaremetalService{client: c} }
func (c *Client) BaremetalAdditionalIP() *BaremetalAdditionalIPService {
	return &BaremetalAdditionalIPService{client: c}
}
func (c *Client) BaremetalElasticStorage() *BaremetalElasticStorageService {
	return &BaremetalElasticStorageService{client: c}
}
func (c *Client) GPU() *GPUService                     { return &GPUService{client: c} }
func (c *Client) ObjectStorage() *ObjectStorageService { return &ObjectStorageService{client: c} }

func (c *Client) do(ctx context.Context, method, path, contentType string, body []byte) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 4; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(1<<(attempt-1)) * 500 * time.Millisecond):
			}
		}
		raw, err := c.once(ctx, method, path, contentType, body)
		if err == nil {
			return raw, nil
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || !retryable(apiErr.StatusCode) {
			return nil, err
		}
		lastErr = err
	}
	return nil, lastErr
}

func (c *Client) once(ctx context.Context, method, path, contentType string, body []byte) ([]byte, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-token", c.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		code, msg := parseErrorBody(raw)
		return nil, &APIError{StatusCode: resp.StatusCode, Code: code, Message: msg, Body: string(raw)}
	}
	return raw, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any) (json.RawMessage, error) {
	var buf []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = b
	}
	raw, err := c.do(ctx, method, path, "application/json", buf)
	if err != nil {
		return nil, err
	}
	return unwrapJSON(raw)
}

// unwrapJSON pisahin payload dari envelope `{success, code, data}`.
// Kalo body valid JSON tapi tanpa key data, balikin body utuh (list endpoint
// kadang kirim bare array/object). Body kosong tetap error.
func unwrapJSON(raw []byte) (json.RawMessage, error) {
	var env struct {
