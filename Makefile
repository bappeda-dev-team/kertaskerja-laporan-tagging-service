APP_NAME=laporan-tagging-service

.PHONY: all build run myenv clean

# DEFAULT TARGET
all: build

build: $(APP_NAME)

$(APP_NAME): *.go
	@echo ">>> Building $(APP_NAME)..."
	@go build -o $(APP_NAME) .
	@echo ">>> SUCCESS..."

run: build myenv
	@echo ">>> Running $(APP_NAME)..."
	./$(APP_NAME)

myenv:
	@echo "REQUIRED ENV"
	@echo "PERENCANAAN_DB_URL: $(PERENCANAAN_DB_URL)"

clean:
	@echo "CLEANING UP"
	rm -f $(APP_NAME)
