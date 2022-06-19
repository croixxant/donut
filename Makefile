NAME=donut
COVERAGE_OUT=cover.out
COVERAGE_HTML=cover.html

.PHONY: build
build:
	go build -o bin/$(NAME)

.PHONY: test
test:
	go test -cover ./... -coverprofile=$(COVERAGE_OUT)

.PHONY: cover
cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

.PHONY: clean
clean:
	rm -f bin/$(NAME)
	rm -f $(COVERAGE_OUT)
	rm -f $(COVERAGE_HTML)