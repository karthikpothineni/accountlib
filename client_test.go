package accountlib

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"accountlib/httprequest"
)

var (
	accountData = map[string][]byte{
		"7eb322ba-57f6-465c-b600-79f26ac7fdc3": []byte(`{"data": {"id":"7eb322ba-57f6-465c-b600-79f26ac7fdc3"}}`),
		"cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb": []byte(`{"data": {"id":"cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb"}`),
	}
)

// ClientTestSuite - test suite for account client
type ClientTestSuite struct {
	suite.Suite
	client *Client
}

// requestHandlerMock - mocks the request handler
type requestHandlerMock struct{}

// SetupTest - This will run before all the tests in the suite
func (s *ClientTestSuite) SetupSuite() {
	options := &ClientOptions{}
	s.client = NewClient(options)
	s.client.handler = &requestHandlerMock{}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

// MakeRequest - function for mocking client MakeRequest
func (r *requestHandlerMock) MakeRequest(specs *httprequest.RequestSpecifications) (statusCode int, body []byte, headers http.Header, err error) {
	if specs.HTTPMethod == http.MethodGet {
		return r.handleGetRequests(specs.URL)
	} else if specs.HTTPMethod == http.MethodPost {
		return r.handlePostRequests(specs.Params)
	} else if specs.HTTPMethod == http.MethodDelete {
		return r.handleDeleteRequests(specs.URL)
	}
	return 0, nil, nil, errors.New("invalid http method")
}

// handleGetRequests - helper for handling mocked GET requests
func (r *requestHandlerMock) handleGetRequests(url string) (statusCode int, body []byte, headers http.Header, err error) {
	urlInfo := strings.Split(url, "/")
	accountID := urlInfo[len(urlInfo)-1]
	if responseBody, ok := accountData[accountID]; ok {
		return http.StatusOK, responseBody, nil, nil
	}
	return http.StatusNotFound, nil, nil, nil
}

// handlePostRequests - helper for handling mocked POST requests
func (r *requestHandlerMock) handlePostRequests(params []byte) (statusCode int, body []byte, headers http.Header, err error) {
	for key, val := range accountData {
		if strings.Contains(string(params), key) {
			return http.StatusCreated, val, nil, nil
		}
	}
	return http.StatusConflict, nil, nil, nil
}

// handleDeleteRequests - helper for handling mocked DELETE requests
func (r *requestHandlerMock) handleDeleteRequests(url string) (statusCode int, body []byte, headers http.Header, err error) {
	urlInfo := strings.Split(url, "/")
	urlSuffix := urlInfo[len(urlInfo)-1]
	urlSuffixInfo := strings.Split(urlSuffix, "?")
	accountID := urlSuffixInfo[0]
	if _, ok := accountData[accountID]; ok {
		return http.StatusNoContent, nil, nil, nil
	}
	return http.StatusNotFound, nil, nil, nil
}

// TestNewClientWithOptions - tests account client object creation with options
func TestNewClientWithOptions(t *testing.T) {
	check := assert.New(t)
	httpClient := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
	options := &ClientOptions{
		HTTPClient: httpClient,
	}
	// create new account client
	client := NewClient(options)
	check.Equal(client.handler, &httprequest.RequestHandler{
		HTTPClient: httpClient,
	})
}

// TestNewClientWithOutOptions - tests account client object creation without options
func TestNewClientWithOutOptions(t *testing.T) {
	check := assert.New(t)
	httpClient := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
	// create new account client
	client := NewClient(nil)
	check.NotEqual(client.handler, &httprequest.RequestHandler{
		HTTPClient: httpClient,
	})
}

// TestFetchAccountSuccessStatusCode - tests an account fetch with successful status code
func (s *ClientTestSuite) TestFetchAccountSuccessStatusCode() {
	check := assert.New(s.T())
	accountID := "7eb322ba-57f6-465c-b600-79f26ac7fdc3"

	// fetch account
	accountData, _ := s.client.Fetch(accountID)
	check.Equal(accountData.ID, accountID)
}

// TestFetchAccountFailureStatusCode - tests an account fetch with failure status code
func (s *ClientTestSuite) TestFetchAccountFailureStatusCode() {
	check := assert.New(s.T())
	incorrectAccountID := "57f6-465c"

	// fetch account
	accountData, err := s.client.Fetch(incorrectAccountID)
	check.Equal(accountData, (*AccountData)(nil))
	check.Contains(err.Error(), "resource not found")
}

// TestFetchAccountEmptyAccount - tests an account fetch with empty account id
func (s *ClientTestSuite) TestFetchAccountEmptyAccount() {
	check := assert.New(s.T())

	// fetch account
	accountData, err := s.client.Fetch("")
	check.Equal(accountData, (*AccountData)(nil))
	check.Contains(err.Error(), "invalid account id")
}

// TestFetchAccountInvalidResponse - tests an account fetch with invalid response
func (s *ClientTestSuite) TestFetchAccountInvalidResponse() {
	check := assert.New(s.T())
	accountID := "cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb"

	// fetch account
	accountData, err := s.client.Fetch(accountID)
	check.Equal(accountData, (*AccountData)(nil))
	check.Contains(err.Error(), "received invalid response")
}

// TestCreateAccountSuccessStatusCode - tests an account creation with successful status code
func (s *ClientTestSuite) TestCreateAccountSuccessStatusCode() {
	check := assert.New(s.T())
	accountID := "7eb322ba-57f6-465c-b600-79f26ac7fdc3"
	orgID := "35eedc2c-0318-40dc-a090-d6f42e7b2754"

	// create account
	accountData, _ := s.client.Create(AccountCreateParams{
		ID:             accountID,
		OrganisationID: orgID,
	})
	check.Equal(accountData.ID, accountID)
}

// TestCreateAccountFailureStatusCode - tests an account creation with failure status code
func (s *ClientTestSuite) TestCreateAccountFailureStatusCode() {
	check := assert.New(s.T())
	conflictAccountID := "57f6-465c"
	orgID := "35eedc2c-0318-40dc-a090-d6f42e7b2754"

	// create account
	accountData, err := s.client.Create(AccountCreateParams{
		ID:             conflictAccountID,
		OrganisationID: orgID,
	})
	check.Equal(accountData, (*AccountData)(nil))
	check.Contains(err.Error(), "request conflict")
}

// TestCreateAccountInvalidResponse - tests an account creation with invalid response
func (s *ClientTestSuite) TestCreateAccountInvalidResponse() {
	check := assert.New(s.T())
	accountID := "cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb"
	orgID := "35eedc2c-0318-40dc-a090-d6f42e7b2754"

	// create account
	accountData, err := s.client.Create(AccountCreateParams{
		ID:             accountID,
		OrganisationID: orgID,
	})
	check.Equal(accountData, (*AccountData)(nil))
	check.Contains(err.Error(), "resource created, but received invalid response")
}

// TestDeleteAccountSuccessStatusCode - tests an account deletion with successful status code
func (s *ClientTestSuite) TestDeleteAccountSuccessStatusCode() {
	check := assert.New(s.T())
	accountID := "7eb322ba-57f6-465c-b600-79f26ac7fdc3"
	version := int64(0)

	// delete account
	err := s.client.Delete(accountID, &version)
	check.Equal(err, nil)
}

// TestDeleteAccountFailureStatusCode - tests an account deletion with failure status code
func (s *ClientTestSuite) TestDeleteAccountFailureStatusCode() {
	check := assert.New(s.T())
	incorrectAccountID := "57f6-465c"
	version := int64(0)

	// delete account
	err := s.client.Delete(incorrectAccountID, &version)
	check.Contains(err.Error(), "resource not found")
}

// TestDeleteAccountEmptyAccount - tests an account deletion with empty account id
func (s *ClientTestSuite) TestDeleteAccountEmptyAccount() {
	check := assert.New(s.T())
	version := int64(0)

	// delete account
	err := s.client.Delete("", &version)
	check.Contains(err.Error(), "invalid account id")
}

// TestDeleteAccountNilVersion - tests an account deletion with nil version
func (s *ClientTestSuite) TestDeleteAccountNilVersion() {
	check := assert.New(s.T())
	accountID := "cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb"

	// delete account
	err := s.client.Delete(accountID, nil)
	check.Contains(err.Error(), "invalid version")
}
