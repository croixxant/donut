VERSION=0.0.1

.PHONY: build
build:
	go build -ldflags "-X main.version=${VERSION}"

.PHONY: test
test:
	go test -cover ./... -coverprofile=cover.out && \
	go tool cover -html=cover.out -o cover.html
