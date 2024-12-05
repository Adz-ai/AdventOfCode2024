# Configuration
ROOT_DIR := $(shell pwd)
TOOLS_DIR := $(ROOT_DIR)/tools
LINT_BIN := $(TOOLS_DIR)/golangci-lint
SESSION_TOKEN := $(AOC_SESSION_TOKEN)
YEAR := 2024

# Find all day directories once
DAY_DIRS := $(shell find . -type d -name "day*" | sort)

# Default target
.PHONY: all test lint run benchmark fetch-input help
all: test lint fetch-input run

# Prevent make from removing intermediate files
.SECONDARY:

run:
	@echo "##### Running all main.go files in each day folder... #####"
	@for dir in $(DAY_DIRS); do \
		if [ -f "$$dir/main.go" ]; then \
			echo "Running $$dir/main.go..."; \
			(cd "$$dir" && go run main.go) || exit $$?; \
		else \
			echo "No main.go found in $$dir, skipping..."; \
		fi; \
	done

test:
	@echo "##### Running all tests in each day folder... #####"
	@for dir in $(DAY_DIRS); do \
		if [ -f "$$dir/main_test.go" ]; then \
			echo "Running $$dir/main_test.go..."; \
			(cd "$$dir" && go test -v) || exit $$?; \
		else \
			echo "No main_test.go found in $$dir, skipping..."; \
		fi; \
	done

lint: $(LINT_BIN)
	@echo "##### Linting all Go files in each day folder using golangci-lint... #####"
	@for dir in $(DAY_DIRS); do \
		if [ -n "$$(ls "$$dir"/*.go 2>/dev/null)" ]; then \
			echo "Linting Go files in $$dir..."; \
			(cd "$$dir" && $(LINT_BIN) run --config $(ROOT_DIR)/.golangci.yml .) || exit $$?; \
		else \
			echo "No Go files found in $$dir, skipping..."; \
		fi; \
	done

$(LINT_BIN):
	@echo "##### Installing golangci-lint... #####"
	@mkdir -p $(TOOLS_DIR)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(TOOLS_DIR)
	@echo "golangci-lint installed at $(LINT_BIN)"

benchmark:
	@echo "##### Measuring performance of each solution for each day... #####"
	@for dir in $(DAY_DIRS); do \
		if [ -f "$$dir/main.go" ]; then \
			echo "Benchmarking $$dir/main.go..."; \
			(cd "$$dir" && env time -p go run main.go) 2>&1 || exit $$?; \
		else \
			echo "No main.go found in $$dir, skipping..."; \
		fi; \
	done

fetch-input:
	@echo "##### Fetching input for each day folder... #####"
	@if [ -z "$(SESSION_TOKEN)" ]; then \
		echo "Error: AOC_SESSION_TOKEN is not set. Please export it in your environment."; \
		exit 1; \
	fi
	@for dir in $(DAY_DIRS); do \
		DAY=$$(echo $$dir | grep -o '[0-9]*'); \
		if [ -n "$$DAY" ]; then \
			if [ -f "$$dir/input.txt" ]; then \
				echo "Input for Day $$DAY already exists. Skipping..."; \
			else \
				URL="https://adventofcode.com/$(YEAR)/day/$$DAY/input"; \
				echo "Fetching input for Day $$DAY from $$URL..."; \
				curl -sSfL --cookie "session=$(SESSION_TOKEN)" "$$URL" -o "$$dir/input.txt" || \
					{ echo "Failed to fetch input for Day $$DAY from $$URL."; exit 1; }; \
			fi; \
		else \
			echo "Invalid day folder: $$dir. Skipping..."; \
		fi; \
	done

help:
	@echo "Available targets:"
	@echo "  help          Display this help message"
	@echo "  all           Run all main.go and main_test.go files in each day folder"
	@echo "  run           Run all main.go files in each day folder"
	@echo "  test          Run all tests in each day folder"
	@echo "  lint          Run golangci-lint in each day folder"
	@echo "  benchmark     Measure performance of each solution for each day"
	@echo "  fetch-input   Fetch puzzle inputs for each day folder"