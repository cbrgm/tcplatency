GO := GO111MODULE=on CGO_ENABLED=0 go

PACKAGES = $(shell go list ./... | grep -v /vendor/)

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

ifndef SHA
	SHA := $(shell git rev-parse --short HEAD)
endif

ifndef DRONE_TAG
	DRONE_TAG := unknown
endif

LDFLAGS += -X main.Version=$(DRONE_TAG)
LDFLAGS += -X main.Revision=$(SHA)
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

.PHONY: check-vendor
check-vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	git update-index --refresh
	git diff-index --quiet HEAD

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: test
test:
	@for PKG in $(PACKAGES); do $(GO) test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

.PHONY: build
build: cmd/tcplatency

.PHONY: cmd/tcplatency
cmd/tcplatency:
	$(GO) build -v -ldflags '-w $(LDFLAGS)' -o ./cmd/tcplatency ./cmd

release:
	@which gox > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	CGO_ENABLED=0 gox -arch="386 amd64 arm" -verbose -ldflags '-w $(LDFLAGS)' -output="dist/tcplatency-${DRONE_TAG}-{{.OS}}-{{.Arch}}" ./cmd