PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64
CMDS := review-create review-query review-export review-checklists compare-create compare-query compare-export

.PHONY: build-all clean

build-all:
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		for cmd in $(CMDS); do \
			echo "Building $$cmd for $$os/$$arch..."; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build \
				-ldflags="-s -w" \
				-o bin/$$os-$$arch/$$cmd$$ext \
				./cmd/$$cmd/; \
		done \
	done
	@echo "✅ contract-atoms build complete!"

clean:
	rm -rf bin/
