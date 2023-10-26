DIST := dist
EXECUTABLE := policy-man

GO ?= go
GOFMT ?= gofmt -l -d -w -s
GOLINT ?= golangci-lint run

TARGETS ?= linux darwin windows
ARCHS ?= amd64
GOFILES := $(shell find . -name "*.go" -type f)
LDFLAGS ?= -X main.commit=$(COMMIT)
EXTLDFLAGS ?=

SWAG := $(shell go env -json | grep 'GOPATH' | cut -d'"' -f4)/bin/swag

.PHONY: all
all: build

.PHONY: install
install: $(GOFILES)
	$(GO) install -mod=readonly -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'
	@echo "\n==>\033[32m Installed policy-man to ${GOPATH}/bin/policy-man\033[m"

.PHONY: build
build: $(EXECUTABLE)

.PHONY: $(EXECUTABLE)
$(EXECUTABLE): $(GOFILES)
	$(GO) build -mod=readonly -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS) -s -w $(LDFLAGS)' -o $@

.PHONY: test
test:
	@$(GO) test -v -cover -coverprofile coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

fmt:
	$(GOFMT) ./

lint:
	$(GOLINT) ./

swag_install:
	$(GO) get -u github.com/swaggo/swag/cmd/swag;go install github.com/swaggo/swag/cmd/swag

api_doc: swag_install
	$(SWAG) fmt;$(SWAG) init --parseDependency --parseInternal -o ./openapi/docs

clean:
	$(GO) clean -modcache -x -i ./...
	find . -name coverage.txt -delete
	find . -name *.tar.gz -delete
	find . -name *.db -delete
	-rm -rf release dist .cover

version:
	@echo $(VERSION)
