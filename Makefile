.PHONY: build clean test install lint deps white verify verify-all analyze

BINARY_NAME=brandkit
BUILD_DIR=bin
CMD_DIR=./cmd/svg

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v ./...

install:
	go install $(CMD_DIR)

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run ./...

# Generate white icons from orig: remove background, convert to white, center, verify
white: build
	@for orig in brands/*/icon_orig.svg; do \
		dir=$$(dirname $$orig); \
		brand=$$(basename $$dir); \
		echo "Processing $$brand..."; \
		$(BUILD_DIR)/$(BINARY_NAME) white $$orig -o $$dir/icon_white.svg; \
	done

# Verify SVG files in a single directory (non-recursive)
verify: build
	$(BUILD_DIR)/$(BINARY_NAME) verify brands/

# Verify all SVG files recursively (for CI)
verify-all: build
	$(BUILD_DIR)/$(BINARY_NAME) verify-all brands/

# Analyze all SVG files for centering
analyze: build
	$(BUILD_DIR)/$(BINARY_NAME) analyze brands/ --fix

build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
