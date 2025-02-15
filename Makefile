# Makefile for building and distributing the neatfile application

APP_NAME = neatfile
VERSION = 1.0.0
BUILD_DIR = build
DIST_DIR = dist

PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all clean build dist brew

all: clean build dist

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)

build:
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} go build -o $(BUILD_DIR)/$(APP_NAME)-$${platform%/*}-$${platform#*/} main.go; \
	done

dist:
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		zip -j $(DIST_DIR)/$(APP_NAME)-$(VERSION)-$${platform%/*}-$${platform#*/}.zip $(BUILD_DIR)/$(APP_NAME)-$${platform%/*}-$${platform#*/}; \
	done

install: build
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

brew:
	@mkdir -p $(DIST_DIR)
	@echo "class Neatfile < Formula" > $(DIST_DIR)/neatfile.rb
	@echo "  desc \"NeatFile is a tool to clean up files by removing comments and empty lines\"" >> $(DIST_DIR)/neatfile.rb
	@echo "  homepage \"https://github.com/AKSarav/neatfile\"" >> $(DIST_DIR)/neatfile.rb
	@echo "  url \"https://github.com/AKSarav/neatfile/releases/download/v$(VERSION)/neatfile-$(VERSION)-darwin-amd64.zip\"" >> $(DIST_DIR)/neatfile.rb
	@echo "  version \"$(VERSION)\"" >> $(DIST_DIR)/neatfile.rb
	@echo "  sha256 \"$(shell shasum -a 256 $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64.zip | awk '{print $$1}')\"" >> $(DIST_DIR)/neatfile.rb
	@echo "  def install" >> $(DIST_DIR)/neatfile.rb
	@echo "    bin.install \"$(APP_NAME)\"" >> $(DIST_DIR)/neatfile.rb
	@echo "  end" >> $(DIST_DIR)/neatfile.rb
	@echo "end" >> $(DIST_DIR)/neatfile.rb