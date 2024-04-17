APP_BASE_NAME	= $(BUILD_DIR)/couture
BUILD_DIR 		= build
COVERAGE_OUT	= $(BUILD_DIR)/coverage.out
TEST_OUT		= $(BUILD_DIR)/test_results.xml
ENVIRONMENT		?= developer-$(USER)
AWS_REGION		?= us-west-2
GOPRIVATE		= github.com/gaggle-net/*
VERSION			= $(shell cat VERSION)

.PHONY: check build
check: cls clean tidy test

rebuild: clean build
build:
	@goreleaser build --clean

.PHONY: cls
cls:
	@clear

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)/

.PHONY: test
test:
	@mkdir -p $(BUILD_DIR)/
	go mod verify

	@command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
	golangci-lint run ./... | tee $(BUILD_DIR)/lint_results.xml
	@echo

	@command -v gotestsum > /dev/null || go install gotest.tools/gotestsum@latest
	gotestsum --junitfile=$(TEST_OUT) --packages=./... -- --tags=$(tags) --coverprofile=$(COVERAGE_OUT) -covermode=atomic -race

	@command -v gocov > /dev/null || go install github.com/axw/gocov/gocov@v1.0.0
	@gocov convert $(COVERAGE_OUT) | gocov-xml > $(BUILD_DIR)/coverage.xml

.PHONY: tidy
tidy:
	@go mod tidy
	@go fix ./...

.PHONY: release
release: build
	git tag $(VERSION)
	git push origin $(VERSION)
