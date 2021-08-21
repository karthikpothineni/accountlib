package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"

	"accountlib"
)

// initClient - initializes the account client
func initClient() *accountlib.Client {
	options := &accountlib.ClientOptions{}
	client := accountlib.NewClient(options)
	return client
}

// createAccount - example function for creating an account
func createAccount(client *accountlib.Client) (*accountlib.AccountData, error) {
	country := "GB"

	// create account
	accountData, err := client.Create(accountlib.AccountCreateParams{
		ID:             uuid.New().String(),
		OrganisationID: "cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb",
		Type:           "accounts",
		Attributes: &accountlib.AccountCreateAttributes{
			Country:                &country,
			BaseCurrency:           "GBP",
			BankID:                 "400300",
			BankIDCode:             "GBDSC",
			Bic:                    "NWBKGB22",
			ProcessingService:      "ABC Bank",
			UserDefinedInformation: "Some important info",
			ValidationType:         "card",
			ReferenceMask:          "############",
			AcceptanceQualifier:    "same_day",
			Name:                   []string{"Samantha Holder"},
			AlternativeNames:       []string{"Sam Holder"},
		},
	},
	)

	return accountData, err
}

// fetchAccount - example function for fetching an account based on account id
func fetchAccount(client *accountlib.Client, accountID string) (*accountlib.AccountData, error) {
	accountData, err := client.Fetch(accountID)
	return accountData, err
}

// deleteAccount - example function for deleting an account based on account id, version
func deleteAccount(client *accountlib.Client, accountID string, version *int64) error {
	err := client.Delete(accountID, version)
	return err
}

func main() {

	// init account client
	client := initClient()

	// create account
	createResponse, err := createAccount(client)
	if err != nil {
		fmt.Printf("Error occurred while creating account - %s\n", err.Error())
		os.Exit(1)
	}
	if accountJson, err := json.Marshal(createResponse); err == nil {
		fmt.Printf("Account created successfully.\nResponse: %s\n", string(accountJson))
	} else {
		fmt.Println("Invalid response received")
		os.Exit(1)
	}

	// fetch account
	fetchResponse, err := fetchAccount(client, createResponse.ID)
	if err != nil {
		fmt.Printf("Error occurred while fetching account - %s\n", err.Error())
		os.Exit(1)
	}
	if accountJson, err := json.Marshal(fetchResponse); err == nil {
		fmt.Printf("Account fetched successfully.\nResponse: %s\n", string(accountJson))
	} else {
		fmt.Println("Invalid response received")
		os.Exit(1)
	}

	// delete account
	err = deleteAccount(client, createResponse.ID, createResponse.Version)
	if err != nil {
		fmt.Println("Unable to delete account")
		os.Exit(1)
	}
	fmt.Println("Account deleted successfully")
}
