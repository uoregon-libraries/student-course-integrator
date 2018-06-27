.PHONY: all binaries clean

GO=vgo
GOFMT=gofmt -s -l -w

all: binaries

binaries:
	$(GO) build -o bin/sci-server github.com/uoregon-libraries/student-course-integrator

format:
	@$(GOFMT) main.go
	@find src -name "*.go" | xargs $(GOFMT)

clean:
	rm bin/* -f
