SUBDIRS := contract-atoms contract-review contract-comparison
SKILLS_DIR := ../skills
DIST_DIR := ../dist
VERSION ?= 1.0.0
NATIVE_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
NATIVE_ARCH := $(shell uname -m)
ifeq ($(NATIVE_ARCH),x86_64)
  NATIVE_ARCH := amd64
endif
ifeq ($(NATIVE_ARCH),aarch64)
  NATIVE_ARCH := arm64
endif
PLATFORM := $(NATIVE_OS)-$(NATIVE_ARCH)

.PHONY: build build-all clean test release release-native package dist dist-clean install

# 只编译当前平台
build:
	@echo "=== Preparing foundation ==="
	@cd foundation && go mod tidy
	@for dir in $(SUBDIRS); do \
		echo ""; \
		echo "=== Building $$dir ($(PLATFORM)) ==="; \
		cd $$dir && go mod tidy && $(MAKE) build && cd ..; \
	done
	@echo ""
	@echo "🎉 Build complete ($(PLATFORM))!"

build-all:
	@echo "=== Preparing foundation ==="
	@cd foundation && go mod tidy
	@for dir in $(SUBDIRS); do \
		echo ""; \
		echo "=== Building $$dir ==="; \
		cd $$dir && go mod tidy && $(MAKE) build-all && cd ..; \
	done
	@echo ""
	@echo "🎉 All builds complete!"

clean:
	@for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done

dist-clean:
	@rm -rf $(DIST_DIR)
	@echo "🧹 Cleaned $(DIST_DIR)/"

test:
	@echo "=== Running foundation tests ==="
	@cd foundation && go test ./... -v

# 构建并打包到 skills/ 目录（SKILL.md + scripts/run.sh + run.ps1 + 跨平台二进制）
release: build-all
	@echo ""
	@echo "=== Packaging to skills/ ==="
	@for dir in $(SUBDIRS); do \
		echo "  → Packaging $$dir..."; \
		mkdir -p $(SKILLS_DIR)/$$dir/scripts/bin; \
		cp $$dir/SKILL.md $(SKILLS_DIR)/$$dir/; \
		cp $$dir/scripts/run.sh $(SKILLS_DIR)/$$dir/scripts/; \
		chmod +x $(SKILLS_DIR)/$$dir/scripts/run.sh; \
		if [ -f $$dir/scripts/run.ps1 ]; then \
			cp $$dir/scripts/run.ps1 $(SKILLS_DIR)/$$dir/scripts/; \
		fi; \
		rm -rf $(SKILLS_DIR)/$$dir/scripts/bin; \
		mkdir -p $(SKILLS_DIR)/$$dir/scripts/bin; \
		cp -r $$dir/bin/* $(SKILLS_DIR)/$$dir/scripts/bin/; \
		find $(SKILLS_DIR)/$$dir/scripts/bin -type f -exec chmod +x {} \;; \
	done
	@echo ""
	@echo "📦 Release packaged to $(SKILLS_DIR)/"
	@echo ""
	@echo "Release contents:"
	@for dir in $(SUBDIRS); do \
		count=$$(find $(SKILLS_DIR)/$$dir/scripts/bin -type f 2>/dev/null | wc -l | tr -d ' '); \
		echo "  $$dir: SKILL.md + run.sh + $$count binaries"; \
	done

# 只编译当前平台并打包到 skills/ 目录（本地开发用）
release-native: build
	@echo ""
	@echo "=== Packaging to skills/ ($(PLATFORM) only) ==="
	@for dir in $(SUBDIRS); do \
		echo "  → Packaging $$dir..."; \
		mkdir -p $(SKILLS_DIR)/$$dir/scripts/bin/$(PLATFORM); \
		cp $$dir/SKILL.md $(SKILLS_DIR)/$$dir/; \
		cp $$dir/scripts/run.sh $(SKILLS_DIR)/$$dir/scripts/; \
		chmod +x $(SKILLS_DIR)/$$dir/scripts/run.sh; \
		if [ -f $$dir/scripts/run.ps1 ]; then \
			cp $$dir/scripts/run.ps1 $(SKILLS_DIR)/$$dir/scripts/; \
		fi; \
		cp $$dir/bin/$(PLATFORM)/* $(SKILLS_DIR)/$$dir/scripts/bin/$(PLATFORM)/; \
		find $(SKILLS_DIR)/$$dir/scripts/bin/$(PLATFORM) -type f -exec chmod +x {} \;; \
	done
	@echo ""
	@echo "📦 Release packaged to $(SKILLS_DIR)/ ($(PLATFORM) only)"

# 编译当前平台 + 安装到 AI 编码工具 skills 目录
# 用法: make install [TOOL=codebuddy|claude|opencode] [TARGET=project|personal]
TOOL ?= codebuddy
TARGET ?= project
install: release-native
	@bash install.sh --tool=$(TOOL) --target=$(TARGET)

# 按平台拆分打包为 tsign-skills-{version}-{platform}.tar.gz
# 用法: make package VERSION=1.0.0
package: release
	@bash scripts/package.sh --version=$(VERSION) --output=$(DIST_DIR)

# dist 是 package 的别名
dist: package
