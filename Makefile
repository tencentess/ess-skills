# 顶层 Makefile — 委托给 toolkit/
VERSION ?= 1.0.0
TOOL ?= codebuddy
TARGET ?= project

.PHONY: build build-all clean test release release-native package dist dist-clean install

# 只编译当前平台
build:
	@$(MAKE) -C toolkit build

build-all:
	@$(MAKE) -C toolkit build-all

clean:
	@$(MAKE) -C toolkit clean

dist-clean:
	@$(MAKE) -C toolkit dist-clean

test:
	@$(MAKE) -C toolkit test

release:
	@$(MAKE) -C toolkit release VERSION=$(VERSION)

release-native:
	@$(MAKE) -C toolkit release-native VERSION=$(VERSION)

package:
	@$(MAKE) -C toolkit package VERSION=$(VERSION)

dist:
	@$(MAKE) -C toolkit dist VERSION=$(VERSION)

# 编译当前平台并安装
# 用法: make install [TOOL=codebuddy|claude|opencode] [TARGET=project|personal]
install:
	@$(MAKE) -C toolkit install TOOL=$(TOOL) TARGET=$(TARGET)
