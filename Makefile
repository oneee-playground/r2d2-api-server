.PHONY: lint
lint:
	@golangci-lint run

UNIT_TEST_PKG = $(shell go list ./... | grep -v /test)

.PHONY: unit-test
unit-test:
	go test -v -race $(UNIT_TEST_PKG)

COVER_OUT?="unit.out"

.PHONY: unit-test-cov
unit-test-cov:
	go test -v -race $(UNIT_TEST_PKG) -coverprofile=$(COVER_OUT) -covermode=atomic

.PHONY: gen-models
gen-models:
	@go generate ./internal/infra/data/ent/model/...
	@go mod tidy

MOCK_PKG = $(shell go list ./... | grep -v /internal/infra/data/ent/model)

.PHONY: gen-mocks
gen-mocks:
	@go generate $(MOCK_PKG)