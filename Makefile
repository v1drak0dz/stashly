APP_NAME := stashly
CMD_DIR := ./cmd/$(APP_NAME)
BUILD_DIR := ./bin

PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: all build clean $(PLATFORMS)

all: build

# Regra genérica para build
define build_target
$1/$2:
	@echo "Building for $1 ($2)..."
	@GOOS=$1 GOARCH=$2 go build -o $(BUILD_DIR)/$(APP_NAME)_$1_$2$(if $(filter $1,windows),.exe,) $(CMD_DIR)
endef

$(foreach platform,$(PLATFORMS),$(eval $(call build_target,$(word 1,$(subst /, ,$(platform))),$(word 2,$(subst /, ,$(platform))))))

# Agrupadores
build-linux: linux/amd64 linux/arm64
build-macos: darwin/amd64 darwin/arm64
build-windows: windows/amd64 windows/arm64

build: $(PLATFORMS)

run: build
	@echo Running $(APP_NAME)...
	@$(BUILD_DIR)/$(APP_NAME).exe

clean:
	@echo Cleaning...
	@rm -rf $(BUILD_DIR)

lint:
	@echo Linting...
	@golangci-lint run ./...

test:
	@echo Testing...
	@go test ./... -v
