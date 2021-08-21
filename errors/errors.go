package accounterrors

import (
	"fmt"
	"net/http"
)

// ErrorMap - holds error message for respective status code
var ErrorMap = map[int]string{
	http.StatusBadRequest:          "bad request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
	http.StatusNotFound:            "resource not found",
	http.StatusMethodNotAllowed:    "incorrect http method",
	http.StatusNotAcceptable:       "incorrect content type",
	http.StatusConflict:            "request conflict",
	http.StatusTooManyRequests:     "too many requests",
	http.StatusInternalServerError: "internal server error",
	http.StatusBadGateway:          "bad gateway",
	http.StatusServiceUnavailable:  "service unavailable",
	http.StatusGatewayTimeout:      "gateway timeout",
}

// handleErrorStatusCode - returns an error based on the status code
func HandleErrorStatusCode(statusCode int, response []byte) (err error) {
	if errMsg, ok := ErrorMap[statusCode]; ok {
		err = fmt.Errorf("%s: %s", errMsg, string(response))
	} else {
		err = fmt.Errorf("internal error: %s", string(response))
	}
	return err
}
