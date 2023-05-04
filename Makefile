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
run-shortener: 
	$(GOCMD) run ./shortener/cmd/v1/main.go http
run-user:
	$(GOCMD) run ./user/cmd/v1/main.go http
	
.PHONY: deps run-auth run-shortener run-user

