AccountLib
==========
AccountLib is a client library for interacting with account API's

## Description
This application is responsible for providing Create, Fetch, Delete functionalities for account API. Internally this library uses in-built net/http library for making http requests.

## Setup
### Local
1. Clone the repository under GOPATH
2. Install dependencies using ```go mod download```

### Docker
Docker compose internally runs the linter, tests before building the application. If there is any error with linter or tests, build will be failed. Run docker compose using 

```docker-compose up --build```
### Run Linter Without Docker
```golangci-lint run -v -c golangci.yml```
### Run Tests Without Docker
```go test -v -cover ./...```

## Example
### Prerequisite:
1. Clone [interview-accountapi](https://github.com/form3tech-oss/interview-accountapi) repository
2. Run ```docker-compose up```

### Execution
1. Run the example using ```go run examples/account.go```

### Output
```$xslt
Account created successfully.
Response: {"attributes":{"alternative_names":["Sam Holder"],"bank_id":"400300","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB22","country":"GB","name":["Samantha Holder"]},"id":"97246f53-cb03-45d2-9479-d01dab2071c1","organisation_id":"cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb","type":"accounts","version":0}
Account fetched successfully.
Response: {"attributes":{"alternative_names":["Sam Holder"],"bank_id":"400300","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB22","country":"GB","name":["Samantha Holder"]},"id":"97246f53-cb03-45d2-9479-d01dab2071c1","organisation_id":"cca3d6ba-cdb1-11eb-be5c-bfc51b0459bb","type":"accounts","version":0}
Account deleted successfully
```

## Code Coverage
Current code coverage is more than **90%**


