PROJECTNAME=$(shell basename "$(PWD)")

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

compile:
	@-$(MAKE) -s go-compile

release:
	@-$(MAKE) -s go-compile-release

test:
	@-$(MAKE) -s go-test

run:
	@-$(MAKE) -s go-compile go-run

clean:
	@-$(MAKE) go-clean


go-compile: go-clean go-build

go-compile-release: go-clean go-build-release

go-run:
	@echo "  >  Run..."
	./bookmarks.api

go-test:
	@echo "  >  Go test..."
	go test -race -v ./...

go-build:
	@echo "  >  Building binary..."
	go build -o bookmarks.api ./main.go

go-build-release:
	@echo "  >  Building binary..."
	GOOS=linux go build -race -ldflags="-s -w" -tags prod -o bookmarks.api ./main.go

go-clean:
	@echo "  >  Cleaning build cache"
	go clean ./...
	rm -f ./bookmarks.api

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
