# Golang Flags
GOPATH ?= $(GOPATH:):./vendor
GOFLAGS ?= $(GOFLAGS:)
GO=go

# Tests
TESTS=.

TESTFLAGS ?= -short
TESTFLAGSEE ?= -test.short

# Packages lists
TE_PACKAGES=$(shell go list ./... | grep -v vendor)
TE_PACKAGES_COMMA=$(shell echo $(TE_PACKAGES) | tr ' ' ',')


all: dist

dist: test-te test-te-race

do-cover-file:
	@echo "mode: count" > cover.out

test-te: do-cover-file
	@echo Testing TE


	@echo "Packages to test: "$(TE_PACKAGES)

	@for package in $(TE_PACKAGES); do \
		echo "Testing "$$package; \
		$(GO) test $(GOFLAGS) -run=$(TESTS) $(TESTFLAGS) -test.v -test.timeout=2000s -covermode=count -coverprofile=cprofile.out -coverpkg=$(ALL_PACKAGES_COMMA) $$package || exit 1; \
		if [ -f cprofile.out ]; then \
			tail -n +2 cprofile.out >> cover.out; \
			rm cprofile.out; \
		fi; \
	done

test-te-race:
	@echo Testing TE race conditions

	@echo "Packages to test: "$(TE_PACKAGES)

	@for package in $(TE_PACKAGES); do \
		echo "Testing "$$package; \
		$(GO) test $(GOFLAGS) -race -run=$(TESTS) -test.timeout=4000s $$package || exit 1; \
	done


cover:
	@echo Opening coverage info in browser. If this failed run make test first

	$(GO) tool cover -html=cover.out
	$(GO) tool cover -html=ecover.out
