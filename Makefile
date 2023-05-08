# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

deps:
	$(GOGET) -v ./...
run-auth: 
	$(GOCMD) run ./auth/cmd/v1/main.go http
run-shortener-http: 
	$(GOCMD) run ./shortener/cmd/v1/main.go http
run-shortener-grpc:
	$(GOCMD) run ./shortener/cmd/v1/main.go grpc
run-shortener-consumer:
	$(GOCMD) run ./shortener/cmd/v1/main.go consumer
run-user:
	$(GOCMD) run ./user/cmd/v1/main.go http
	
.PHONY: deps run-auth run-shortener-http run-shortener-grpc run-shortener-consumer run-user

