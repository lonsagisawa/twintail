.PHONY: build build-css install-deps clean dev

build: build-css
	go build -tags prod -o twintail

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
	rm -f twintail
