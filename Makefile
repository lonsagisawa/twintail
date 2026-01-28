.PHONY: build build-css install-deps clean dev

build: build-css build-linux-amd64 build-linux-arm64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -tags prod -ldflags="-s -w" -o twintail-linux-amd64

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -tags prod -ldflags="-s -w" -o twintail-linux-arm64

build-css:
	npm run build-css

watch-css:
	npm run watch-css

dev:
	@trap 'kill 0' EXIT; \
	npm run watch-css & \
	air

install-deps:
	npm install
	go install github.com/air-verse/air@latest

clean:
	rm -f twintail-linux-amd64 twintail-linux-arm64
