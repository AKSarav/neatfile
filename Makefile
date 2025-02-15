# Makefile for building and distributing the neatfile application

APP_NAME = neatfile
VERSION = 0.0.1
BUILD_DIR = build
DIST_DIR = dist
LICENSE = MIT

PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all clean build dist brew

all: clean build dist

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)

build:
	@echo "Building neatfile version $(VERSION)"
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} go build -o $(BUILD_DIR)/$(APP_NAME)-$${platform%/*}-$${platform#*/} main.go; \
	done

dist: build
	@echo "Creating distribution packages"
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		zip -j $(DIST_DIR)/$(APP_NAME)-$(VERSION)-$${platform%/*}-$${platform#*/}.zip $(BUILD_DIR)/$(APP_NAME)-$${platform%/*}-$${platform#*/}; \
	done

install: build
	@echo "Installing neatfile version $(VERSION)"
	@cp $(BUILD_DIR)/$(APP_NAME)-$$(go env GOOS)-$$(go env GOARCH) $(PWD)/$(APP_NAME)
	@echo "Installation complete. Adding $(APP_NAME) to PATH."
	@CURRENT_SHELL=$$(basename $$SHELL); \
	if [ "$$CURRENT_SHELL" = "bash" ]; then \
		echo 'export PATH=$(PWD):$$PATH' >> ~/.bashrc; \
		echo 'Added to ~/.bashrc. Please run `source ~/.bashrc` to update your PATH.'; \
	elif [ "$$CURRENT_SHELL" = "zsh" ]; then \
		echo 'export PATH=$(PWD):$$PATH' >> ~/.zshrc; \
		echo 'Added to ~/.zshrc. Please run `source ~/.zshrc` to update your PATH.'; \
	elif [ "$$CURRENT_SHELL" = "fish" ]; then \
		echo 'set -U fish_user_paths $(PWD) $$fish_user_paths' >> ~/.config/fish/config.fish; \
		echo 'Added to ~/.config/fish/config.fish. Please run `source ~/.config/fish/config.fish` to update your PATH.'; \
	else \
		echo 'Unknown shell. Please add $(PWD) to your PATH manually.'; \
	fi