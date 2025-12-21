.PHONY: build-all clean test run-lab-01 run-lab-02 run-lab-03 run-lab-04 run-lab-05 run-lab-06 run-lab-07 run-lab-08 run-lab-09 run-lab-10

# Build all labs
build-all:
	@echo "Building all labs..."
	@for lab in 01_labs/greedy_heuristics 02_labs/greedy_regret_heuristics 03_labs/local_search 04_labs/local_search_candidate_moves 05_labs/local_search_deltas 06_labs/local_search_extensions 07_labs/large_neighborhood_search 08_labs/global_convexity 09_labs/hybrid_evolutionary_algorithm 10_labs/variable_neighborhood_search; do \
		echo "Building $$lab..."; \
		go build -o $$lab/cmd/main $$lab/cmd/main.go || echo "Failed to build $$lab"; \
	done

# Clean all build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@find . -name "main" -type f -delete
	@find . -name "*.test" -type f -delete
	@find . -name "*.out" -type f -delete
	@rm -rf output/
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test ./pkg/common/... -v

# Run individual labs
run-lab-01:
	@cd 01_labs/greedy_heuristics/cmd && go run main.go

run-lab-02:
	@cd 02_labs/greedy_regret_heuristics/cmd && go run main.go

run-lab-03:
	@cd 03_labs/local_search/cmd && go run main.go

run-lab-04:
	@cd 04_labs/local_search_candidate_moves/cmd && go run main.go

run-lab-05:
	@cd 05_labs/local_search_deltas/cmd && go run main.go

run-lab-06:
	@cd 06_labs/local_search_extensions/cmd && go run main.go

run-lab-07:
	@cd 07_labs/large_neighborhood_search/cmd && go run main.go

run-lab-08:
	@cd 08_labs/global_convexity/cmd && go run main.go

run-lab-09:
	@cd 09_labs/hybrid_evolutionary_algorithm/cmd && go run main.go

run-lab-10:
	@cd 10_labs/variable_neighborhood_search/cmd && go run main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run ./... || echo "golangci-lint not installed, skipping..."

