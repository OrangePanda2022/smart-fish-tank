# Smart Fish Tank — Backend Build Makefile
# Location: /home/lenovo/SHIT2/smart-fish-tank/Makefile

SHELL := /bin/bash
.PHONY: all build clean \
	auth gateway mind sensor tank

GOOS   ?= linux
GOARCH ?= amd64
GO_BUILDER_IMAGE ?= golang:1.25-alpine
HOST_UID := $(shell id -u)
HOST_GID := $(shell id -g)

BACKEND_DIR := $(CURDIR)/Backend
BIN_DIR     := $(CURDIR)/bin

# Optional version injection (harmless if mains don't define these vars)
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)
STATIC_LDFLAGS := $(LDFLAGS) -linkmode external -extldflags=-static
CGO_TAGS := osusergo netgo sqlite_omit_load_extension
DOCKER_GO := docker run --rm -v $(CURDIR):/src -w /src $(GO_BUILDER_IMAGE) /bin/sh -lc
ALPINE_BUILD_DEPS := apk add --no-cache build-base ca-certificates
GO_BIN := /usr/local/go/bin/go
CGO_CFLAGS := -D_LARGEFILE64_SOURCE

# ----------------------------------------------------------------------------
# DEFAULT: build all services on the host
# ----------------------------------------------------------------------------
all: build

build: auth gateway mind sensor tank

# --- auth (CGO/sqlite; build static for Alpine runtime) ---
auth:
	@echo "==> Building auth ..."
	@mkdir -p $(BIN_DIR)/auth
	$(DOCKER_GO) '$(ALPINE_BUILD_DEPS) && cd /src/Backend/auth && CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=gcc CGO_CFLAGS="$(CGO_CFLAGS)" $(GO_BIN) build -tags "$(CGO_TAGS)" -trimpath -ldflags "$(STATIC_LDFLAGS)" -o /src/bin/auth/auth ./cmd/server && chmod +x /src/bin/auth/auth && chown $(HOST_UID):$(HOST_GID) /src/bin/auth/auth'
	@chmod +x $(BIN_DIR)/auth/auth

# --- gateway (pure Go; static binary works on Alpine) ---
gateway:
	@echo "==> Building gateway ..."
	@mkdir -p $(BIN_DIR)/gateway
	cd $(BACKEND_DIR)/gateway && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/gateway/gateway ./cmd/gateway
	@chmod +x $(BIN_DIR)/gateway/gateway

# --- mind (pure Go; static binary works on Alpine) ---
mind:
	@echo "==> Building mind ..."
	@mkdir -p $(BIN_DIR)/mind
	cd $(BACKEND_DIR)/mind && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/mind/mind ./cmd/server
	@chmod +x $(BIN_DIR)/mind/mind

# --- sensor (CGO/sqlite; build static for Alpine runtime) ---
sensor:
	@echo "==> Building sensor ..."
	@mkdir -p $(BIN_DIR)/sensor
	$(DOCKER_GO) '$(ALPINE_BUILD_DEPS) && cd /src/Backend/sensor && CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=gcc CGO_CFLAGS="$(CGO_CFLAGS)" $(GO_BIN) build -tags "$(CGO_TAGS)" -trimpath -ldflags "$(STATIC_LDFLAGS)" -o /src/bin/sensor/sensor ./cmd/server && chmod +x /src/bin/sensor/sensor && chown $(HOST_UID):$(HOST_GID) /src/bin/sensor/sensor'
	@chmod +x $(BIN_DIR)/sensor/sensor

# --- tank (CGO/sqlite; build static for Alpine runtime) ---
tank:
	@echo "==> Building tank ..."
	@mkdir -p $(BIN_DIR)/tank
	$(DOCKER_GO) '$(ALPINE_BUILD_DEPS) && cd /src/Backend/tank && CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=gcc CGO_CFLAGS="$(CGO_CFLAGS)" $(GO_BIN) build -tags "$(CGO_TAGS)" -trimpath -ldflags "$(STATIC_LDFLAGS)" -o /src/bin/tank/tank ./cmd/server && chmod +x /src/bin/tank/tank && chown $(HOST_UID):$(HOST_GID) /src/bin/tank/tank'
	@chmod +x $(BIN_DIR)/tank/tank

# ----------------------------------------------------------------------------
# CLEAN: remove only executables, preserve configs / data / docker-compose.yaml
# ----------------------------------------------------------------------------
clean:
	@echo "==> Removing compiled binaries (preserving configs, data, docker-compose) ..."
	rm -f $(BIN_DIR)/auth/auth
	rm -f $(BIN_DIR)/gateway/gateway
	rm -f $(BIN_DIR)/mind/mind
	rm -f $(BIN_DIR)/sensor/sensor
	rm -f $(BIN_DIR)/tank/tank
	@echo "    Done."
