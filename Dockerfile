FROM golang:1.15-alpine

MAINTAINER Karthik Pothineni

LABEL service=accountlib

# Install golang linter
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.32.0

WORKDIR /etc/accountlib

ENV CGO_ENABLED=0

COPY . .

RUN go mod download

# Run linter
RUN golangci-lint run -v -c golangci.yml

# Run tests
RUN go test -v -cover ./...
