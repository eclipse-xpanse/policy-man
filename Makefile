DIST := dist
EXECUTABLE := policy-man

GO ?= go
GOFMT ?= gofmt -l -d -w -s

TARGETS ?= linux darwin windows
ARCHS ?= amd64
GOFILES := $(shell find . -name "*.go" -type f)
LDFLAGS ?= -X main.version=$(VERSION) -X main.commit=$(COMMIT)
EXTLDFLAGS ?=

.PHONY: all
all: build

.PHONY: install
install: $(GOFILES)
	$(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'
	@echo "\n==>\033[32m Installed policy-man to ${GOPATH}/bin/policy-man\033[m"

.PHONY: build
build: $(EXECUTABLE)

.PHONY: $(EXECUTABLE)
$(EXECUTABLE): $(GOFILES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS) -s -w $(LDFLAGS)' -o $@

.PHONY: test
test:
	@$(GO) test -v -cover -coverprofile coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

fmt:
	$(GOFMT) ./

clean:
	$(GO) clean -modcache -x -i ./...
	find . -name coverage.txt -delete
	find . -name *.tar.gz -delete
	find . -name *.db -delete
	-rm -rf release dist .cover

version:
	@echo $(VERSION)
