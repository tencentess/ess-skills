PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64
NATIVE_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
NATIVE_ARCH := $(shell uname -m)
ifeq ($(NATIVE_ARCH),x86_64)
  NATIVE_ARCH := amd64
endif
ifeq ($(NATIVE_ARCH),aarch64)
  NATIVE_ARCH := arm64
endif

.PHONY: build build-all clean

# 只编译当前平台
build:
	@echo "Building compare-workflow for $(NATIVE_OS)/$(NATIVE_ARCH)..."
	@CGO_ENABLED=0 GOOS=$(NATIVE_OS) GOARCH=$(NATIVE_ARCH) go build \
		-ldflags="-s -w" \
		-o bin/$(NATIVE_OS)-$(NATIVE_ARCH)/compare-workflow \
		./cmd/compare-workflow/
	@echo "✅ contract-comparison build complete ($(NATIVE_OS)/$(NATIVE_ARCH))"

# 编译全部平台
build-all:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		echo "Building compare-workflow for $$os/$$arch..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build \
			-ldflags="-s -w" \
			-o bin/$$os-$$arch/compare-workflow$$ext \
			./cmd/compare-workflow/; \
	done
	@echo "✅ contract-comparison build complete!"

clean:
	rm -rf bin/
