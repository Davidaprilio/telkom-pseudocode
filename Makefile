ifeq ($(OS),Windows_NT)
BIN := telupsc.exe
else
BIN := telupsc
endif

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
FILE := $(if $(ARGS),$(ARGS),./examples/if.telu)
OUT ?=

.PHONY: help tidy build run build-example go-example clean

help:
	@echo "Available targets:"
	@echo "  make tidy                      - Install/update Go dependencies"
	@echo "  make build                     - Build telupsc CLI binary"
	@echo "  make run FILE=path/file.telu   - Run .telu file via source"
	@echo "  make build-example FILE=...    - Build .telu into binary via source"
	@echo "  make go-example FILE=...       - Export transpiled .go via source"
	@echo "  make clean                     - Remove generated CLI binary"

tidy:
	go mod tidy

build:
	go build -o $(BIN) ./cmd/telupsc

run:
	go run ./cmd/telupsc run $(FILE)

build-example:
	go run ./cmd/telupsc build $(FILE) $(if $(OUT),-o $(OUT),)

go-example:
	go run ./cmd/telupsc go $(FILE) $(if $(OUT),-o $(OUT),)

clean:
	rm -f telupsc telupsc.exe
