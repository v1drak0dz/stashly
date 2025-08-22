APP_NAME := stashly
CMD_DIR := ./cmd/$(APP_NAME)
BUILD_DIR := ./bin

.PHONY: all build run clean lint test

all: build

build-linux:
	@echo "Building for linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

build-windows:
	@echo "Building for windows..."
	@go build -o $(BUILD_DIR)/$(APP_NAME).exe $(CMD_DIR)

build: build-linux build-windows

run: build
	@echo Running $(APP_NAME)...
	@$(BUILD_DIR)/$(APP_NAME).exe

clean:
	@echo Cleaning...
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)

lint:
	@echo Linting...
	@golangci-lint run ./...

test:
	@echo Testing...
	@go test ./... -v
