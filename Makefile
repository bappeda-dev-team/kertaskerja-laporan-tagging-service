APP_NAME=laporan-tagging-service

# DEFAULT TARGET
.PHONY: all
all: build

# Build binary
.PHONY: build
build:
	@echo ">>> Building $(APP_NAME)..."
	@go build -o $(APP_NAME) .
	@echo ">>> SUCCESS..."

# Run with env
.PHONY: run
run: build
	@echo ">>> Running $(APP_NAME)..."
	./$(APP_NAME)
