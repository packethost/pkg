help: ## Print this help
	@grep --no-filename -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/:.*##/·/' | sort | column -ts '·' -c 120

all: example ## Build example binary

-include lint.mk
-include rules.mk

test: ## Run tests
	go test -v ./...

coverage: ## Generate coverage report
	CGO_ENABLED=1 go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Enable this later
#verify: gofumpt prettier lint # # Verify code style, is lint free, freshness ...
verify: goimports lint ## Verify code style, is lint free, freshness ...
	git diff | (! grep .)

# Enable this later
#fix: gofumpt-fix prettier-fix ## Fix code formatting errors

tools: ${toolsBins} ## Build Go based build tools

run: run-example ## Run examplelog

.PHONY: all coverage help run test tools verify
