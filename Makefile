.PHONY: all deps binaries test clean dbconf

GO=go
GOFMT=gofmt -s -l -w

all: binaries

deps:
	$(GO) get

binaries:
	@./scripts/hackversion.sh
	$(GO) install ./src/...
	$(GO) build -o bin/sci github.com/uoregon-libraries/student-course-integrator

format:
	@$(GOFMT) main.go
	@find . -name "*.go" | xargs $(GOFMT)

validate:
	./scripts/validate.sh

test:
	@$(GO) test ./... | grep -v "^?.*no test files"

clean:
	rm bin/* -f

dbconf:
	./scripts/makedbconf.sh
