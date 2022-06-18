NAME=donut
VERSION=0.0.1
COVERAGE_OUT=cover.out
COVERAGE_HTML=cover.html

.PHONY: build
build:
	go build -ldflags "-X main.version=${VERSION}" -o bin/$(NAME)

.PHONY: test
test:
	go test -cover ./... -coverprofile=$(COVERAGE_OUT)

.PHONY: cover
cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
