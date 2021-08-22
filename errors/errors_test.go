package accounterrors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestServiceUnavailableError - test if service unavailable error is found for a status code
func TestServiceUnavailableError(t *testing.T) {
	check := assert.New(t)
	err := HandleErrorStatusCode(http.StatusServiceUnavailable, []byte(`{"data": null}`))
	check.Contains(err.Error(), "service unavailable")
}

// TestUnknownError - test if error is not found for a status code
func TestUnknownError(t *testing.T) {
	check := assert.New(t)
	err := HandleErrorStatusCode(http.StatusVariantAlsoNegotiates, []byte(`{"data": null}`))
	check.Contains(err.Error(), "internal error")
}

// TestGateWayTimeoutError - test if gateway timeout error is found for a status code
func TestGateWayTimeoutError(t *testing.T) {
	check := assert.New(t)
	err := HandleErrorStatusCode(http.StatusGatewayTimeout, []byte(`{"data": null}`))
	check.Contains(err.Error(), "gateway timeout")
}

// TestEmptyResponse - test if response is empty
func TestEmptyResponse(t *testing.T) {
	check := assert.New(t)
	err := HandleErrorStatusCode(0, []byte(``))
	check.Contains(err.Error(), "internal error")
}
