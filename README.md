Setup
--

- [Install Go](https://golang.org/dl/)
- Set up your GOPATH: https://golang.org/doc/code.html#GOPATH
  - Consider putting `$GOPATH/bin` into your path
- Put this repository into `$GOPATH/src/github.com/uoregon-libraries/student-course-integrator`
- Get goose for database migrations: `go get -u bitbucket.org/liamstask/goose/...`
- `make`

Settings file
---

Copy example.conf to sci.conf and modify it as needed.  SCI will look for this
file at `/etc/sci.conf`, then `./sci.conf`.  You can pass the `-c` flag to
specify a custom location as well, e.g., `./bin/sci server -c /tmp/dummysettings.conf`.

All settings can be overridden with environment variables prefixed with "SCI_".
In production, use this to avoid storing sensitive values in `sci.conf`:

    export SCI_DB=sciusername:password@tcp(localhost:port)/scidbname
    export SCI_SESSION_SECRET=blah

Set up the database
---

You'll need to set up a database and user for SCI to store its faculty/course
association.  For development, this is trivial by using the supplied
docker-compose configuration:

```bash
  docker-compose up -d
```

That would generate a database with username, password, and database name of "sci".

This project uses [goose](https://bitbucket.org/liamstask/goose) for managing
the database tables.  Create a database with mysql or mariadb, set up a
`db/dbconf.yml` file:

```yaml
development:
  driver: mysql
  open: user:password@tcp(localhost:port)/databasename
```

And finally, run `goose up`.

The Makefile has a target which generates `db/dbconf.yml` for you if your DB
config lives in `/etc/sci.conf` or `./sci.conf`:

```bash
make dbconf
```

This eliminates the need to generate your db config, but note that it will
overwrite any existing `db/dbconf.yml`.

Run the server
---

```bash
./bin/sci server
```

Development
---

As mentioned above, consider using docker to ease development by giving you a
preconfigured database.  You can then populate the "sci" database tables with
any fake (or real) courses and user ids.  To log in as any given user, make
sure you have DEBUG set to true in your configuration (or `export SCI_DEBUG =
1` to temporarily set this up), then visit the app with a "debuguser" query
argument.  For example, `http://localhost:8080/?debuguser=jechols`.

If you install [entr](http://www.entrproject.org/), you can speed up your
development loop by running [`./scripts/devloop.sh`](./scripts/devloop.sh),
which runs [`makerun.sh`](./scripts/makerun.sh) whenever `entr` detects a
change to any file or directory under `src/`.
