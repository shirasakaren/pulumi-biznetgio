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

// unwrapJSON extracts the payload from the `{success, code, data}` envelope.
// If the body is valid JSON but has no data key, it returns the body as-is (list
// endpoints sometimes send a bare array/object). An empty body is still an error.
func unwrapJSON(raw []byte) (json.RawMessage, error) {
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err == nil &&
		len(env.Data) > 0 && string(env.Data) != "null" {
		return env.Data, nil
	}
	if json.Valid(raw) {
		return raw, nil
	}
	return nil, fmt.Errorf("biznetgio: empty response body")
}

// jsonInt64 is tolerant: it accepts a number OR a numeric string on the wire (the Rust SDK
// uses deserialize_number_from_string for keypair_id and friends).
type jsonInt64 int64

func (n *jsonInt64) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*n = jsonInt64(v)
		return nil
	}
	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*n = jsonInt64(v)
	return nil
}

func decodeJSON[T any](raw json.RawMessage) (T, error) {
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("biznetgio: decode response: %w", err)
	}
	return out, nil
}

func withStatus(p, status string) string {
	if status == "" {
		return p
	}
	return p + "?status=" + url.QueryEscape(status)
}

func strField(d map[string]any, key string) string {
	v, ok := d[key]
	if !ok {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	}
	return ""
}

func intField(d map[string]any, key string) int64 {
	v, ok := d[key]
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return int64(x)
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	}
	return 0
}
