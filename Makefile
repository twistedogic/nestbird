.PHONY: build install uninstall clean test

BINARY_NAME=nestbird
INSTALL_PATH=/usr/local/bin
SYSTEMD_PATH=/etc/systemd/system

# Build the binary
build:
	go build -o $(BINARY_NAME) .

# Build for Linux (cross-compile from macOS)
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .

# Install the binary
install: build
	install -m 755 $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

# Install systemd unit
install-service:
	install -m 644 nestbird.service $(SYSTEMD_PATH)/nestbird.service
	systemctl daemon-reload

# Uninstall binary
uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)

# Uninstall systemd unit
uninstall-service:
	rm -f $(SYSTEMD_PATH)/nestbird.service
	systemctl daemon-reload

# Full install (binary + service)
full-install: install install-service
	@echo "Run 'systemctl enable --now nestbird' to start the service"

# Full uninstall
full-uninstall: uninstall uninstall-service

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux-amd64

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
