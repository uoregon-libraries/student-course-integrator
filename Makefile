.PHONY: all binaries clean dbconf

GO=vgo
GOFMT=gofmt -s -l -w

all: binaries

binaries:
	$(GO) build -o bin/sci-server github.com/uoregon-libraries/student-course-integrator
	@# This is helpful to get code completion tools working while vgo is still in transition
	@go install ./src/...

format:
	@$(GOFMT) main.go
	@find src -name "*.go" | xargs $(GOFMT)

clean:
	rm bin/* -f

dbconf:
	./scripts/makedbconf.sh
