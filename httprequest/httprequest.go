package httprequest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

// http request constants
const (
	defaultRetryCount            = 3
	defaultMaxIdleConnection     = 100
	defaultKeepAliveTime         = 30 * time.Second
	defaultIdleConnectionTimeout = 30 * time.Second
	defaultTimeout               = 5 * time.Second
	defaultRequestType           = "application/json"
)

// default transport and retry codes
var (
	defaultTransport = &http.Transport{
		DialContext: (&net.Dialer{
			KeepAlive: defaultKeepAliveTime,
		}).DialContext,
		MaxIdleConns:        defaultMaxIdleConnection,
		IdleConnTimeout:     defaultIdleConnectionTimeout,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConnsPerHost: defaultMaxIdleConnection,
	}
	defaultRetryStatusCodes = []int{
		http.StatusRequestTimeout,
		http.StatusGatewayTimeout,
		http.StatusServiceUnavailable,
	}
)

type RequestHandlerIface interface {
	MakeRequest(specs *RequestSpecifications) (statusCode int, body []byte, headers http.Header, err error)
}

// RequestSpecifications - controls the each http requests behaviour
type RequestSpecifications struct {
	URL        string
	HTTPMethod string
	Params     []byte
	Timeout    int
	RetryCount int
}

// RequestHandler - holds http client
type RequestHandler struct {
	HTTPClient *http.Client
}

// NewRequestHandler  - returns RequestHandler object
func NewRequestHandler(customClient *http.Client) *RequestHandler {
	if customClient != nil {
		return &RequestHandler{
			HTTPClient: customClient,
		}
	}
	httpClient := &http.Client{}
	httpClient.Transport = defaultTransport
	httpClient.Timeout = time.Duration(defaultTimeout)
	return &RequestHandler{
		HTTPClient: httpClient,
	}
}

// MakeRequest - prepares request and makes an API call
func (r *RequestHandler) MakeRequest(specs *RequestSpecifications) (statusCode int, body []byte, headers http.Header, err error) {
	baseBackOffTime := 100 * time.Millisecond
	requestCount := 1

	// prepare request
	newHandler, newRequest, err := r.prepareRequest(specs)
	if err != nil {
		return statusCode, nil, nil, err
	}

	// handle retries using exponential backoff strategy
	for requestCount <= specs.RetryCount {
		// sending the request
		statusCode, body, headers, err = sendRequest(newHandler, newRequest)
		if checkRetryRequired(statusCode) || err != nil {
			time.Sleep(time.Duration(baseBackOffTime))
			baseBackOffTime = 2 * baseBackOffTime
		} else {
			break
		}
		requestCount++
	}

	return
}

// prepareRequest - returns customized request handler with default values if not exclusively specified
func (r *RequestHandler) prepareRequest(specs *RequestSpecifications) (*http.Client, *http.Request, error) {
	//Create request
	req, err := http.NewRequest(specs.HTTPMethod, specs.URL, nil)
	if err != nil {
		err = fmt.Errorf("unable to create http request. error: %s", err.Error())
		return r.HTTPClient, req, err
	}
	//check and set retry count
	if specs.RetryCount == 0 {
		specs.RetryCount = defaultRetryCount
	}
	// check and set request timeout
	if specs.Timeout != 0 {
		r.HTTPClient.Timeout = time.Duration(specs.Timeout) * time.Second
	}
	// add post headers and body
	if specs.HTTPMethod == http.MethodPost {
		// handle POST request
		req.Header.Add("Content-type", defaultRequestType)
		body := prepareRequestBody(specs.Params)
		req.Body = body
	}
	return r.HTTPClient, req, nil
}

// prepareRequestBody - converts []byte to readcloser
func prepareRequestBody(params []byte) io.ReadCloser {
	var body *bytes.Buffer
	body = bytes.NewBuffer(params)
	return ioutil.NopCloser(body)
}

// checkRetryRequired - checks if retry is required based on the status code
func checkRetryRequired(statusCode int) bool {
	retryFlag := false
	for _, retryCode := range defaultRetryStatusCodes {
		retryFlag = retryFlag || (retryCode == statusCode)
	}
	return retryFlag
}

// sendRequest - sends HTTP request
func sendRequest(newHandler *http.Client, newRequest *http.Request) (int, []byte, http.Header, error) {
	// send http request
	resp, err := newHandler.Do(newRequest)
	if err != nil {
		if os.IsTimeout(err) {
			err = fmt.Errorf("timeout encountered. error: %s", err.Error())
			return http.StatusRequestTimeout, nil, nil, err
		}
		err = fmt.Errorf("failed to send request. Error: %s", err.Error())
		return 0, nil, nil, err
	}

	// read response body
	defer resp.Body.Close()
	body := new(bytes.Buffer)
	_, readError := body.ReadFrom(resp.Body)
	if readError != nil {
		err = fmt.Errorf("failed to read response body. error: %s", err.Error())
		return resp.StatusCode, nil, nil, err
	}

	return resp.StatusCode, body.Bytes(), resp.Header, nil
}
