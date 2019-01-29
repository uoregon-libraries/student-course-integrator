.PHONY: all deps binaries test clean dbconf

GO=go
GOFMT=gofmt -s -l -w

all: binaries

deps: version
	$(GO) get

binaries: version
	$(GO) install ./src/...
	$(GO) build -o bin/sci

format: version
	@$(GOFMT) main.go
	@find . -name "*.go" | xargs $(GOFMT)

validate: version
	./scripts/validate.sh

test: version
	@$(GO) test ./... | grep -v "^?.*no test files"

version:
	@go generate ./src/version
	@chmod a+w src/version/build.go 2>/dev/null || true

clean:
	rm bin/* -f
	rm src/version/commit.go -f

dbmigrate:
	./scripts/dbmigrate.sh

ci: deps binaries validate test
