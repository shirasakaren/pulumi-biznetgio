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
