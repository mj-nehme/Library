# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

# ARGS 
RUN_ARGS=$(filter-out $@,$(MAKECMDGOALS))
TEST_ARGS=-v -p 1

target: build 

# Install project dependencies
install-deps:
	go mod download

# Binaries
BINARY_DIR=bin
SRV_BINARY_NAME=server
SRV_BINARY_PATH=$(BINARY_DIR)/$(SRV_BINARY_NAME)

# Main
SRV_MAIN=./main.go

#Integration Test
INTG_TEST_DIR=./tests
INTG_TEST=$(INTG_TEST_DIR)/...
UNIT_TEST=./...

# Build the main application server
build: install-deps
	$(GOBUILD) -o $(SRV_BINARY_PATH) $(SRV_MAIN)

# Run the main application server
run:
	$(GORUN) $(SRV_MAIN) $(RUN_ARGS)

# Unit Test
unit-test:
	$(GOTEST) $(TEST_ARGS) $(UNIT_TEST)

# Integration Test
integration-test:
	$(GOTEST) $(TEST_ARGS) $(INTG_TEST)

test: unit-test integration-test

# Clean up compiled files
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)