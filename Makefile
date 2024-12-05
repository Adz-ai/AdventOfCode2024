ROOT_DIR := $(shell pwd)
TOOLS_DIR := $(ROOT_DIR)/tools
LINT_BIN := $(TOOLS_DIR)/golangci-lint
SESSION_TOKEN := $(AOC_SESSION_TOKEN) # Ensure to export this in your environment

all: test lint fetch-input run

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
			cd $$dir && go test -v || exit $$?; \
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
			cd $$dir && $(LINT_BIN) run --config $(ROOT_DIR)/.golangci.yml . || exit $$?; \
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

benchmark:
	@echo "##### Measuring performance of each solution for each day... #####"
	@for dir in $(shell find . -type d -name "day*" | sort); do \
		if [ -f $$dir/main.go ]; then \
			echo "Benchmarking $$dir/main.go..."; \
			cd $$dir && { env time -p go run main.go || exit $$?; } 2>&1; \
			cd $(ROOT_DIR); \
		else \
			echo "No main.go found in $$dir, skipping..."; \
		fi; \
	done

fetch-input:
	@echo "##### Fetching input for each day folder... #####"
	@if [ -z "$(AOC_SESSION_TOKEN)" ]; then \
		echo "Error: AOC_SESSION_TOKEN is not set. Please export it in your environment."; \
		exit 1; \
	fi
	@YEAR=2024; \
	for dir in $(shell find . -type d -name "day*" | sort); do \
		DAY=$$(echo $$dir | grep -o '[0-9]*'); \
		if [ -n "$$DAY" ]; then \
			if [ -f $$dir/input.txt ]; then \
				echo "Input for Day $$DAY already exists. Skipping..."; \
			else \
				URL="https://adventofcode.com/$$YEAR/day/$$DAY/input"; \
				echo "Fetching input for Day $$DAY from $$URL..."; \
				curl -sSfL --cookie "session=$(AOC_SESSION_TOKEN)" $$URL -o $$dir/input.txt || { \
					echo "Failed to fetch input for Day $$DAY from $$URL."; exit 1; }; \
			fi; \
		else \
			echo "Invalid day folder: $$dir. Skipping..."; \
		fi; \
	done


help:
	@echo "Available targets:"
	@echo "  help    	Display this help message"
	@echo "  all     	Run all main.go and main_test.go files in each day folder"
	@echo "  run     	Run all main.go files in each day folder"
	@echo "  test    	Run all tests in each day folder"
	@echo "  lint    	Run golangci-lint in each day folder"
	@echo "  benchmark 	Measure performance of each solution for each day"
	@echo "  fetch-input Fetch puzzle inputs for each day folder"

