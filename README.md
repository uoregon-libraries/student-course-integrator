Setup
--

- [Install Go](https://golang.org/dl/)
- Grab vgo: `go get -u golang.org/x/vgo`
- Get goose for database migrations: `go get -u bitbucket.org/liamstask/goose/...`
- `make`

If you've installed [entr](http://www.entrproject.org/), you can speed up your
development loop by running [`./scripts/devloop.sh`](./scripts/devloop.sh),
which runs [`makerun.sh`](./scripts/makerun.sh) whenever `entr` detects a
change to any file or directory under `src/`.

**Note**: The make recipe uses vgo, which doesn't install compiled packages in
a location that's friendly for things like `gocode` to give you auto-completion
features.  If you want compiled package files generated, you'll have to
override the `INSTALL` variable when you run `make`:

```bash
export INSTALL=1
make -e # Or ./scripts/devloop.sh, which already adds the -e flag to make
```

Settings file
---

Copy example.conf to sci.conf and modify it as needed.  SCI will look for this
file at `/etc/sci.conf`, then `./sci.conf`.  You can pass the `-c` flag to
specify a custom location as well, e.g., `./bin/sci-server -c /tmp/dummysettings.conf`.

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
./bin/sci-server
```
