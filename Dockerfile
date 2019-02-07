# Build step
FROM golang:1.11 as build-env

RUN apt-get update -y
RUN apt-get install -y upx-ucl

# Get the goose binary built in a portable way
RUN go get -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags="no_sqlite no_psql" -i -o /go/bin/goose ./cmd/goose
RUN upx /go/bin/goose

# Clone the github repo so our first build is cached by docker
RUN go get -d github.com/uoregon-libraries/student-course-integrator
WORKDIR /go/src/github.com/uoregon-libraries/student-course-integrator
RUN make deps && make

# Add local code and rebuild the sci binary
ADD Makefile .
ADD db .
ADD main.go .
ADD src .
ADD static .
ADD templates .
RUN make
RUN upx ./bin/sci

# Production step
FROM alpine:latest

# Dependencies
RUN apk update && apk add ca-certificates && apk add bash && apk add mysql-client && rm -rf /var/cache/apk/*

# These are required for the app to run
RUN mkdir -p /var/sci/csvin
RUN mkdir -p /var/sci/csvout

COPY --from=build-env /go/src/github.com/uoregon-libraries/student-course-integrator/bin/sci /app/sci
COPY --from=build-env /go/bin/goose /usr/local/bin/goose

# Add local files that didn't need to be compiled or otherwise processed
COPY ./static /app/static
COPY ./templates /app/templates
COPY ./db /app/db
COPY ./scripts/dbmigrate.sh /app/db/migrate.sh
COPY ./scripts/docker-entry.sh /entrypoint.sh
COPY ./scripts/wait_for_database /usr/local/bin/wait_for_database

WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
