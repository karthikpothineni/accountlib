package httprequest

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HttpTestSuite - test suite for http request
type HTTPTestSuite struct {
	suite.Suite
	requestHandler        *RequestHandler
	requestSpecifications *RequestSpecifications
	url                   string
}

// SetupTest - This will run before all the tests in the suite
func (s *HTTPTestSuite) SetupSuite() {
	s.requestHandler = NewRequestHandler(nil)
	s.url = "http://localhost:8080/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	s.requestSpecifications = &RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        s.url,
	}
	httpmock.ActivateNonDefault(s.requestHandler.HTTPClient)
}

// BeforeTest - This will run before each test
func (s *HTTPTestSuite) BeforeTest(suiteName, testName string) {
	httpmock.Reset()
}

// AfterTest - This will run after test suite
func (s *HTTPTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}

// TestNewRequestHandler - test request handler object creation
func TestNewRequestHandler(t *testing.T) {
	check := assert.New(t)
	requestHandler := NewRequestHandler(nil)
	check.Equal(requestHandler.HTTPClient.Timeout, defaultTimeout)
	check.NotEqual(requestHandler.HTTPClient.Transport, nil)
}

// TestNewRequestHandlerWithCustomClient - test request handler object creation with custom client
func TestNewRequestHandlerWithCustomClient(t *testing.T) {
	check := assert.New(t)
	customClientTimeout := time.Duration(10) * time.Second
	requestHandler := NewRequestHandler(&http.Client{
		Timeout: customClientTimeout,
	})
	check.Equal(requestHandler.HTTPClient.Timeout, customClientTimeout)
}

// TestMakeRequestSuccessResponse - tests a successful api call
func (s *HTTPTestSuite) TestMakeRequestSuccessResponse() {
	check := assert.New(s.T())
	nullData := `{"data": null}`

	// mock http request
	httpmock.RegisterResponder(http.MethodGet, s.url,
		httpmock.NewStringResponder(http.StatusOK, nullData))

	// make http request
	statusCode, response, _, err := s.requestHandler.MakeRequest(s.requestSpecifications)
	if err == nil {
		check.Equal(statusCode, http.StatusOK)
		check.Equal(string(response), nullData)
	}

	// get the amount of calls for the registered responder
	info := httpmock.GetCallCountInfo()
	check.Equal(1, info[http.MethodGet+" "+s.url])
}

// TestMakeRequestFailureResponse - tests a failure api call
func (s *HTTPTestSuite) TestMakeRequestFailureResponse() {
	check := assert.New(s.T())

	// mock http request
	httpmock.RegisterResponder(http.MethodGet, s.url,
		httpmock.NewStringResponder(http.StatusBadGateway, ``))

	// make http request
	statusCode, response, _, err := s.requestHandler.MakeRequest(s.requestSpecifications)
	if err == nil {
		check.Equal(statusCode, http.StatusBadGateway)
		check.Equal(string(response), ``)
	}

	// get the amount of calls for the registered responder
	info := httpmock.GetCallCountInfo()
	check.Equal(1, info[http.MethodGet+" "+s.url])
}

// TestMakeRequestCreationError - tests http request create error
func (s *HTTPTestSuite) TestMakeRequestCreationError() {
	check := assert.New(s.T())

	// make http request
	_, _, _, err := s.requestHandler.MakeRequest(&RequestSpecifications{
		HTTPMethod: "*?",
	})
	check.Contains(err.Error(), "unable to create http request")
}

// TestMakeRequestWrongHttpMethod - tests with wrong http method
func (s *HTTPTestSuite) TestMakeRequestWrongHttpMethod() {
	check := assert.New(s.T())

	// make http request
	_, _, _, err := s.requestHandler.MakeRequest(&RequestSpecifications{
		HTTPMethod: "TEST",
	})
	check.Contains(err.Error(), "failed to send request")
}

// TestPrepareRequestCustomTimeout - tests prepare request with custom timeout
func (s *HTTPTestSuite) TestPrepareRequestCustomTimeout() {
	check := assert.New(s.T())
	customTimeout := 10

	// make http request
	_, _, _ = s.requestHandler.prepareRequest(&RequestSpecifications{
		HTTPMethod: http.MethodPost,
		Params:     []byte(""),
		Timeout:    customTimeout,
	})
	check.Equal(s.requestHandler.HTTPClient.Timeout, time.Duration(customTimeout)*time.Second)
}

// TestRetryRequired - tests a successful retry check
func TestRetryRequired(t *testing.T) {
	check := assert.New(t)
	retryRequired := checkRetryRequired(http.StatusServiceUnavailable)
	check.Equal(retryRequired, true)
}

// TestRetryNotRequired - tests a failure retry check
func TestRetryNotRequired(t *testing.T) {
	check := assert.New(t)
	retryRequired := checkRetryRequired(http.StatusConflict)
	check.Equal(retryRequired, false)
}
