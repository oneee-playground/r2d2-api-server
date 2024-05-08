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