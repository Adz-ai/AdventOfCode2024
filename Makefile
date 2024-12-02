ROOT_DIR := $(shell pwd)
TOOLS_DIR := $(ROOT_DIR)/tools
LINT_BIN := $(TOOLS_DIR)/golangci-lint

all: test lint run

run:
	@echo "##### Running all main.go files in each day folder... #####"
	@for dir in $(shell find . -type d -name "day*" | sort); do \
		if [ -f $$dir/main.go ]; then \
			echo "Running $$dir/main.go..."; \
			cd $$dir && go run main.go; \
			cd $(ROOT_DIR); \
		else \
			echo "No main.go found in $$dir, skipping..."; \
		fi; \
	done

test:
	@echo "##### Running all tests in each day folder... #####"
	@for dir in $(shell find . -type d -name "day*" | sort); do \
		if [ -f $$dir/main_test.go ]; then \
			echo "Running $$dir/main_test.go..."; \
			cd $$dir && go test -v; \
			cd $(ROOT_DIR); \
		else \
			echo "No main_test.go found in $$dir, skipping..."; \
		fi; \
	done

lint: $(LINT_BIN)
	@echo "##### Linting all Go files in each day folder using golangci-lint... #####"
	@for dir in $(shell find . -type d -name "day*" | sort); do \
		if [ -n "$$(ls $$dir/*.go 2>/dev/null)" ]; then \
			echo "Linting Go files in $$dir..."; \
			cd $$dir && $(LINT_BIN) run --config $(ROOT_DIR)/.golangci.yml .; \
			cd $(ROOT_DIR); \
		else \
			echo "No Go files found in $$dir, skipping..."; \
		fi; \
	done



$(LINT_BIN):
	@echo "##### Installing golangci-lint... #####"
	@mkdir -p $(TOOLS_DIR)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR)
	@echo "golangci-lint installed at $(LINT_BIN)"

# Help target
help:
	@echo "Available targets:"
	@echo "  help    Display this help message"
	@echo "  all     Run all main.go and main_test.go files in each day folder"
	@echo "  run     Run all main.go files in each day folder"
	@echo "  test    Run all tests in each day folder"
	@echo "  lint    Run golangci-lint in each day folder"
	@echo "  $(LINT_BIN) Retrieve and install golangci-lint in tools directory"
