# viber00t Makefile

BINARY_NAME=viber00t
INSTALL_DIR=$(HOME)/.local/bin

.PHONY: all build install clean help

all: build

build:
	@echo "Building viber00t..."
	@go build -o $(BINARY_NAME)

install: build
	@echo "Installing viber00t to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed to $(INSTALL_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@$(BINARY_NAME) clean 2>/dev/null || true

uninstall:
	@echo "Uninstalling viber00t..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Uninstalled"

help:
	@echo "viber00t Makefile targets:"
	@echo "  make         - Build viber00t"
	@echo "  make install - Build and install to ~/.local/bin"
	@echo "  make clean   - Remove binary and clean images"
	@echo "  make uninstall - Remove from ~/.local/bin"