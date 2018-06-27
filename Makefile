.PHONY: all binaries clean

GO=vgo

all: binaries

binaries:
	$(GO) build -o bin/sci-server github.com/uoregon-libraries/student-course-integrator

clean:
	rm bin/* -f
