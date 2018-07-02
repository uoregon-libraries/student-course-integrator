.PHONY: all binaries test clean dbconf

GO=vgo
GOFMT=gofmt -s -l -w
INSTALL=0

all: binaries

binaries:
	$(GO) build -o bin/sci-server github.com/uoregon-libraries/student-course-integrator
	@# This is helpful to get code completion tools working while vgo is still in transition
	@[ "$(INSTALL)" -ne "1" ] || ./scripts/install.sh

format:
	@$(GOFMT) main.go
	@find src -name "*.go" | xargs $(GOFMT)

validate:
	find ./src -type f -name "*.go" | sed -s "s|/\w\+\.go||" | sort | uniq | xargs go vet
	find ./src -type f -name "*.go" | sed -s "s|/\w\+\.go||" | sort | uniq | xargs golint
	find ./src -type f -name "*.go" | xargs gofmt -l -s

test:
	@$(GO) test ./... | grep -v "^?.*no test files"

clean:
	rm bin/* -f

dbconf:
	./scripts/makedbconf.sh
