export GO111MODULE = on

REPOSITORY = go.mercari.io/go-emv-code
PACKAGES = $(shell go list ./...)

GO_TEST ?= go test
GO_TEST_TARGET ?= .

REVIEWDOG_ARG ?= -diff="git diff master"

LINT_TOOLS=$(shell cat tools/tools.go | egrep '^\s_ '  | awk '{ print $$2 }')

.PHONY: all
all: test

.PHONY: bootstrap-lint-tools
bootstrap-lint-tools:
	@echo "Installing/Updating tools (dir: $(GOBIN), tools: $(LINT_TOOLS))"
	@go install -tags tools -mod=readonly $(LINT_TOOLS)

.PHONY: test
test:  ## Run go test
	${GO_TEST} -v -race -mod=readonly -run=$(GO_TEST_TARGET) $(PACKAGES)

.PHONY: coverage
coverage:  ## Collect test coverage
	${GO_TEST} -v -race -mod=readonly -run=$(GO_TEST_TARGET) -covermode=atomic -coverpkg=${REPOSITORY}/... -coverprofile=$@.out $(PACKAGES)

.PHONY: reviewdog
reviewdog: bootstrap-lint-tools  ## Run reviewdog
	reviewdog -conf=.reviewdog.yml $(REVIEWDOG_ARG)

.PHONY: help
help:  ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[33m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z\/_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
