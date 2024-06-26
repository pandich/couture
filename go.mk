BUILD_DIR 		= build
COVERAGE_OUT	= $(BUILD_DIR)/coverage.out
TEST_OUT		= $(BUILD_DIR)/test_results.xml
ENVIRONMENT		?= developer-$(USER)
AWS_REGION		?= us-west-2
GOPRIVATE		= github.com/gaggle-net/*
VERSION			= $(shell cat VERSION)

.PHONY: build
build: clean quality test

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)/

.PHONY: test
test: quality
	@mkdir -p $(BUILD_DIR)/
	go mod verify

	@command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./... | tee $(BUILD_DIR)/lint_results.xml
	@echo

	@command -v gotestsum > /dev/null || go install gotest.tools/gotestsum@latest
	gotestsum --junitfile=$(TEST_OUT) --packages=./... -- --tags=$(tags) --coverprofile=$(COVERAGE_OUT) -covermode=atomic -race

	@command -v gocov > /dev/null || go install github.com/axw/gocov/gocov@v1.0.0
	@gocov convert $(COVERAGE_OUT) | gocov-xml > $(BUILD_DIR)/coverage.xml

.PHONY: quality quality-mod quality-lint quality-staticcheck
quality: quality-mod quality-lint quality-staticcheck
quality-mod:
	go mod tidy
	go mod verify
quality-lint:
	@command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint --timeout=10s run ./...
quality-staticcheck:
	@command -v staticcheck > /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -tests -f stylish ./...

.PHONY: tidy
tidy:
	go mod tidy
	go fix ./...


.PHONY: release
release:
	git tag $(VERSION)
	git push origin $(VERSION)
	rm -rf dist/ build/
	goreleaser


.PHONY: docs
docs:
	@echo http://127.0.0.1:20080/pkg/github.com/gagglepanda/couture/
	@godoc -index -play -http 127.0.0.1:20080

cloc:
	@cloc --not-match-f='_test.go$$' --not-match-d='^examples$$' --include-lang=Go .
