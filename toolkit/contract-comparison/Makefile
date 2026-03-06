PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: build-all clean

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
