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
