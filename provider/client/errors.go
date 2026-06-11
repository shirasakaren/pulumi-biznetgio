package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Code       int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("biznetgio: status %d", e.StatusCode)
	}
	return fmt.Sprintf("biznetgio: status %d: %s", e.StatusCode, e.Message)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

func parseErrorBody(body []byte) (int, string) {
	var m struct {
		Code    int    `json:"code"`
		Detail  any    `json:"detail"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(body, &m); err != nil {
		return 0, ""
	}
	msg := string(body)
	if d, ok := m.Detail.(string); ok {
		msg = d
	} else if m.Message != "" {
		msg = m.Message
	} else if m.Error != "" {
		msg = m.Error
	}
	return m.Code, msg
}
