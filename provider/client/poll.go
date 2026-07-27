package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func retryable(status int) bool {
	switch status {
	case http.StatusTooManyRequests, http.StatusInternalServerError,
		http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

// WaitForStatus polls get until status(t) matches a terminal value.
// interval: how often to poll. If ctx has no deadline, a 30m timeout is applied.
func WaitForStatus[T any](
	ctx context.Context,
	interval time.Duration,
	get func(ctx context.Context) (T, error),
	status func(T) string,
	ready, failed []string,
) (T, error) {
	var zero T
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
