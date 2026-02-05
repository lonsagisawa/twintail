.PHONY: build build-frontend install-deps clean dev test test-coverage install uninstall

ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
	BINARY := twintail-linux-amd64
else ifeq ($(ARCH),aarch64)
	BINARY := twintail-linux-arm64
else
	BINARY := twintail-linux-amd64
endif

PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin
SYSTEMD_DIR := /etc/systemd/system

build: build-frontend build-linux-amd64 build-linux-arm64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -tags prod -ldflags="-s -w" -o twintail-linux-amd64 ./cmd/server

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -tags prod -ldflags="-s -w" -o twintail-linux-arm64 ./cmd/server

build-frontend:
	npm run build

dev:
	@trap 'kill 0' EXIT; \
	npm run dev & \
	air

install-deps:
	npm install
	go install github.com/air-verse/air@latest

clean:
	rm -f twintail-linux-amd64 twintail-linux-arm64
	rm -rf static/dist

test:
	go test -v ./...

test-coverage:
	go test -cover ./...

install: $(BINARY)
	install -Dm755 $(BINARY) $(DESTDIR)$(BINDIR)/twintail
	install -Dm644 twintail.service $(DESTDIR)$(SYSTEMD_DIR)/twintail.service
	@echo "Installed twintail to $(DESTDIR)$(BINDIR)/twintail"
	@echo "Installed systemd unit to $(DESTDIR)$(SYSTEMD_DIR)/twintail.service"
	@echo "Run 'systemctl daemon-reload && systemctl enable --now twintail' to start"

uninstall:
	-systemctl stop twintail 2>/dev/null || true
	-systemctl disable twintail 2>/dev/null || true
	rm -f $(DESTDIR)$(BINDIR)/twintail
	rm -f $(DESTDIR)$(SYSTEMD_DIR)/twintail.service
	systemctl daemon-reload
	@echo "Uninstalled twintail"
