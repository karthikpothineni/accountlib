package accountlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"accountlib/errors"
	"accountlib/httprequest"
)

// account api constants
const (
	accountBaseURL = "http://localhost:8080"
	accountPath    = "v1/organisation/accounts"
)

// Client - holds account client information
type Client struct {
	handler httprequest.RequestHandlerIface
}

// ClientOptions - options passed while creating a new client
// Users can control connection pooling by passing a custom http client
type ClientOptions struct {
	HTTPClient *http.Client
}

// AccountCreateParams - holds fields for account creation
// This struct is similar to AccountResponse but excludes some unnecessary fields for creation
type AccountCreateParams struct {
	Attributes     *AccountCreateAttributes `json:"attributes,omitempty"`
	ID             string                   `json:"id,omitempty"`
	OrganisationID string                   `json:"organisation_id,omitempty"`
	Type           string                   `json:"type,omitempty"`
}

// AccountCreateAttributes - holds account attributes for account creation
// This struct is similar to AccountResponseAttributes but excludes some unnecessary fields for creation
type AccountCreateAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
	ProcessingService       string   `json:"processing_service,omitempty"`
	UserDefinedInformation  string   `json:"user_defined_information,omitempty"`
	ValidationType          string   `json:"validation_type,omitempty"`
	ReferenceMask           string   `json:"reference_mask,omitempty"`
	AcceptanceQualifier     string   `json:"acceptance_qualifier,omitempty"`
}

// AccountData - holds complete account response
type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

// AccountAttributes - holds account attribute response
type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	StatusReason            string   `json:"status_reason,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
	ProcessingService       string   `json:"processing_service,omitempty"`
	UserDefinedInformation  string   `json:"user_defined_information,omitempty"`
	ValidationType          string   `json:"validation_type,omitempty"`
	ReferenceMask           string   `json:"reference_mask,omitempty"`
	AcceptanceQualifier     string   `json:"acceptance_qualifier,omitempty"`
}

// NewClient - creates a new account client
func NewClient(options *ClientOptions) (client *Client) {
	client = &Client{}

	// prepare http client
	client.handler = httprequest.NewRequestHandler(options.HTTPClient)

	return client
}

// Fetch - returns the account details based on account id
func (client *Client) Fetch(accountID string) (accountData *AccountData, err error) {
	// validate account id
	if accountID == "" {
		err = errors.New("invalid account id")
		return
	}

	// prepare request specifications
	url := fmt.Sprintf("%s/%s/%s", accountBaseURL, accountPath, accountID)
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
	}

	// make request
	statusCode, response, _, err := client.handler.MakeRequest(requestSpecifications)
	if err != nil {
		return
	}

	// handle status code, response
	if statusCode == http.StatusOK {
		dataResponse := make(map[string]AccountData)
		err = json.Unmarshal(response, &dataResponse)
		if err != nil {
			err = fmt.Errorf("received invalid response. error: %s", err.Error())
			return
		}
		if accountData, ok := dataResponse["data"]; ok {
			return &accountData, nil
		}
	} else {
		err = accounterrors.HandleErrorStatusCode(statusCode, response)
	}

	return
}

// Create - creates an account based on create params
func (client *Client) Create(createParams AccountCreateParams) (accountData *AccountData, err error) {
	// marshal create params
	dataMap := make(map[string]AccountCreateParams)
	dataMap["data"] = createParams
	params, err := json.Marshal(dataMap)
	if err != nil {
		err = fmt.Errorf("unable to marshal create params, error: %s", err.Error())
		return
	}

	// prepare request specifications
	url := fmt.Sprintf("%s/%s", accountBaseURL, accountPath)
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodPost,
		URL:        url,
		Params:     params,
	}

	// make request
	statusCode, response, _, err := client.handler.MakeRequest(requestSpecifications)
	if err != nil {
		return
	}

	// handle status code, response
	if statusCode == http.StatusCreated {
		dataResponse := make(map[string]AccountData)
		err = json.Unmarshal(response, &dataResponse)
		if err != nil {
			err = fmt.Errorf("resource created, but received invalid response. error: %s", err.Error())
			return
		}
		if accountData, ok := dataResponse["data"]; ok {
			return &accountData, nil
		}
	} else {
		err = accounterrors.HandleErrorStatusCode(statusCode, response)
	}

	return
}

// Delete  - deletes an account based on account id and version
func (client *Client) Delete(accountID string, version *int64) (err error) {
	// validate account id, version
	if accountID == "" {
		err = errors.New("invalid account id")
		return
	}
	if version == nil {
		err = errors.New("invalid version")
		return
	}

	// prepare request specifications
	url := fmt.Sprintf("%s/%s/%s?version=%d", accountBaseURL, accountPath, accountID, *version)
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodDelete,
		URL:        url,
	}

	// make request
	statusCode, response, _, err := client.handler.MakeRequest(requestSpecifications)
	if err != nil {
		return
	}

	// handle status code, response
	if statusCode != http.StatusNoContent {
		err = accounterrors.HandleErrorStatusCode(statusCode, response)
	}

	return
}
