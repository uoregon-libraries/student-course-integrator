.PHONY: all binaries test clean dbconf

GO=vgo
GOFMT=gofmt -s -l -w
INSTALL=0

all: binaries

binaries:
	$(GO) build -o bin/sci github.com/uoregon-libraries/student-course-integrator
	@# This is helpful to get code completion tools working while vgo is still in transition
	@[ "$(INSTALL)" -ne "1" ] || ./scripts/install.sh

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
