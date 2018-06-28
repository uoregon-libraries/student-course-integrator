.PHONY: all binaries clean

GO=vgo
GOFMT=gofmt -s -l -w
INSTALL=0

all: binaries

binaries:
	$(GO) build -o bin/sci-server github.com/uoregon-libraries/student-course-integrator
	@# This is helpful to get code completion tools working while vgo is still in transition
	@[[ $(INSTALL) != 1 ]] || go install ./src/...

format:
	@$(GOFMT) main.go
	@find src -name "*.go" | xargs $(GOFMT)

clean:
	rm bin/* -f
