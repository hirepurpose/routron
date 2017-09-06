
# the product we're building
NAME := routron
# the product's main package
MAIN := ./src/main
# fix our gopath
GOPATH := $(GOPATH):$(PWD)

# build and packaging
TARGETS := $(PWD)/bin
PRODUCT := $(TARGETS)/$(NAME)
VERSION ?= $(shell git log --pretty=format:'%h' -n 1)

# build and install
PREFIX ?= /usr/local

# sources
SRC = $(shell find src -name \*.go -print)

# tests
TEST_PKGS = routron

.PHONY: all build install test clean

all: build

build: $(PRODUCT) ## Build the product

$(PRODUCT): $(SRC)
	go build -o $@ $(MAIN)

install: build ## Build and install
	install -m 0755 $(PRODUCT) $(PREFIX)/bin/

test: ## Run tests
	go test -test.v $(TEST_PKGS)

clean: ## Delete the built product and any generated files
	rm -rf $(TARGETS)
