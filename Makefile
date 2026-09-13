GO      ?= go
PKG     := ./cmd/web
PORT    ?= 8080
BIN_DIR := bin
BIN     := $(BIN_DIR)/learngo

ifeq ($(OS),Windows_NT)
BIN := $(BIN).exe
endif

.DEFAULT_GOAL := help

.PHONY: help run build tidy fmt vet test check clean

help:
	$(info LearnGO commands)
	$(info   make run    Run the web app on PORT=$(PORT))
	$(info   make build  Build $(BIN))
	$(info   make tidy   Download and tidy Go modules)
	$(info   make fmt    Format Go source files)
	$(info   make vet    Run go vet)
	$(info   make test   Run tests)
	$(info   make check  Run fmt, vet, and test)
	$(info   make clean  Remove build output)
	@exit 0

run:
	PORT=$(PORT) $(GO) run $(PKG)

build:
	$(GO) build -o $(BIN) $(PKG)

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

check: fmt vet test

clean:
	$(GO) clean
ifeq ($(OS),Windows_NT)
	-if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)
else
	rm -rf $(BIN_DIR)
endif
