PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: build-all clean

build-all:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		echo "Building review-workflow for $$os/$$arch..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build \
			-ldflags="-s -w" \
			-o bin/$$os-$$arch/review-workflow$$ext \
			./cmd/review-workflow/; \
	done
	@echo "✅ contract-review build complete!"

clean:
	rm -rf bin/
