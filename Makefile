.PHONY: all deps binaries test clean dbconf

GO=go
GOFMT=gofmt -s -l -w

all: binaries

deps: version
	$(GO) get

binaries: version
	$(GO) install ./src/...
	$(GO) build -o bin/sci github.com/uoregon-libraries/student-course-integrator

format: version
	@$(GOFMT) main.go
	@find . -name "*.go" | xargs $(GOFMT)

validate: version
	./scripts/validate.sh

test: version
	@$(GO) test ./... | grep -v "^?.*no test files"

version:
	@./scripts/hackversion.sh

clean:
	rm bin/* -f

dbconf:
	./scripts/makedbconf.sh
