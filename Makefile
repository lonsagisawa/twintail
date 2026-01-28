.PHONY: build build-frontend install-deps clean dev

build: build-frontend build-linux-amd64 build-linux-arm64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -tags prod -ldflags="-s -w" -o twintail-linux-amd64

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -tags prod -ldflags="-s -w" -o twintail-linux-arm64

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
