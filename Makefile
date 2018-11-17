include .env

PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
GOFILES=$(wildcard *.go)

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=/tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## compile: Compile the binary.
compile:
	@-$(MAKE) -s go-compile

release:
	@-$(MAKE) -s go-compile-release

test:
	@-$(MAKE) -s go-test

run:
	@-$(MAKE) -s go-compile go-run

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-$(MAKE) go-clean

go-compile: go-clean go-build

go-compile-release: go-clean go-build-release

go-run:
	@echo "  >  Run..."
	./bookmarks.api

go-test:
	@echo "  >  Go test..."
	grc go test -v ./...

go-build:
	@echo "  >  Building binary..."
	go build -v  --toolexec="/usr/bin/time -v" -o bookmarks.api ./*.go

go-build-release:
	@echo "  >  Building binary..."
	GOOS=linux go build ---toolexec="/usr/bin/time -v" ldflags="-s -w" -o bookmarks.api ./*.go

go-clean:
	@echo "  >  Cleaning build cache"
	go clean

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
