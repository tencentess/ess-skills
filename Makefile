# 顶层 Makefile — 委托给 toolkit/
VERSION ?= 1.0.0

.PHONY: build-all clean test release package dist dist-clean

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

package:
	@$(MAKE) -C toolkit package VERSION=$(VERSION)

dist:
	@$(MAKE) -C toolkit dist VERSION=$(VERSION)
