.PHONY: all binaries clean

GO=vgo

all: binaries

binaries:
	$(GO) build -o bin/sci github.com/uoregon-libraries/student-course-integrator/src/cmd/sci

clean:
	rm bin/* -f
