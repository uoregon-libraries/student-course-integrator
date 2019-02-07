This README is not meant to help end-users.  That's a guide we also need but
don't have time to work on right now.

**Note**: For best results, you should compile and test locally.  This is
easier and faster than doing it through docker.

If you opt not to compile and test locally, you'll have to rebuild the images
every time you change code.  This makes the edit-compile-test loop go from
speedy-quick (typical Go) to awful (though still not as bad as Java and Ruby).
This is not officially supported due to the high risk of perpetuating the
belief that slow restarts are acceptable to web developers.

If you want to set things up manually for a fully local setup, the
docker-compose.yml and Dockerfile will help guide you, but the documentation no
longer supports doing this.  Having two sets of documents simply didn't work
out well (especially since the old manual setup still involved docker for the
database).  And while full-local dev is very fast and ever-so-convenient, we
needed a single unified Docker setup for our staging server.  Bleah.

Preliminary Setup
---

### Install and set up Go

Download Go for your OS: https://golang.org/dl/

Just get the latest binary distribution, e.g.:

```bash
cd ~
curl -L "https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz" | tar -xz
```

You'll have a fancy new dir, `~/go`, containing the compiler and tools.

I always move the dir to `~/.go` so my home dir seems less cluttered.

### Set up paths

I just have this in my Linux profile:

```bash
# Remember, I moved it to $HOME/.go - just point this to wherever your Go
# download ended up
export GOROOT=$HOME/.go

# This puts "global" projects under $HOME/projects/go/... - choose a
# destination you like, as this can be almost anywhere
export GOPATH=$HOME/projects/go

# Make sure the Go binaries are in your path as well as any Go tools you
# download using "go get", such as goose (see below)
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

This tells Go where it lives, where you want to download projects, etc.

Verify you set things up nicely by typing `go version`.  You should see
something like "go version go1.11.4 linux/amd64".

Read more about the [GOPATH](https://golang.org/doc/code.html#GOPATH).  If
something isn't making sense, knowing how GOPATH works can help.

### Download SCI

Diehards would recommend using `go get`, but it's sort of magic.  It's just a
lot clearer if you do a git clone.

```bash
mkdir -p $GOPATH/src/github.com/uoregon-libraries
git clone git@github.com:uoregon-libraries/student-course-integrator.git \
          $GOPATH/src/github.com/uoregon-libraries/student-course-integrator
cd $GOPATH/src/github.com/uoregon-libraries/student-course-integrator
```

### Get golint

`golint` is a tool for analyzing your Go code to look for various best
practices, and is required by the test suite:

```bash
go get -u golang.org/x/lint/golint
```

Docker Setup
---

- Build and run all tests so all the necessary generated files are owned by you: `make ci`
- [Install docker-compose](https://docs.docker.com/compose/install/#install-compose)
- Copy `env-example` to `.env` and adjust the configurations as necessary
  - All configuration values in `env-example` will need to be changed
  - Explanations of all configuration options are commented in `./example.conf`.  *You should read these*.
- Copy `./docker-compose.override-example.yml` to `./docker-compose.override-example.yml`
  - This should work without alterations, but you should still look it over and
    make sure it makes sense for your workstation

Development
---

Run the server: `docker-compose up`

When code changes, rebuild and restart the "web" service:

```bash
make
docker-compose restart web
```

You can build the binary by simply running `make`, but it's advised that you
also validate the code and run tests before pushing anything up:

```bash
# Validate, test, and build in a single command
make ci

# Or run the validation, test, and build commands individually
make validate
make test
make
```

SCI's Settings
---

SCI's configuration is very flexible.  You can put config in `/etc/sci.conf` or
`./sci.conf`, or you can pass the `-c` flag to specify a custom location:

```bash
./bin/sci server -c /tmp/dummysettings.conf
```

Additionally, all settings can be overridden with environment variables
prefixed with "SCI_".  In production, use this to avoid storing sensitive
values in `sci.conf`:

```bash
export SCI_DB=sci:p@ssw0rD4lyfe@tcp(maria.mysite.org:3306)/scidb
export SCI_SESSION_SECRET=blah
export SCI_LDAP_SERVER="ldap.mysite.org:389"
export SCI_LDAP_BIND_USER="root"
export SCI_LDAP_BIND_PASS="s3cur3!"
export SCI_LDAP_BASE_DN="dc=ad,dc=mysite,dc=org"
```

The mix of configuration file and environment is how the docker setup works,
letting us control some settings, like the database connection, while putting
others in the devs' hands.

Log in
---

Visit `http://localhost:8080` (assuming you didn't adjust the default port in
the docker override) and you should see the app's "not authorized" page, and
you'll be logged in as "dummyuser" in place of a real authorization.

In debug mode, you can fake a login as any other user by visiting the page with
a "debuguser" query argument.  For example, `http://localhost:8080/?debuguser=jechols`.
The page will have a large, visible warning if it is in debug mode to avoid
accidentally pushing debug to production.

In development, docker auto-loads some seed data.  Use the "debuguser" argument
to sign in as "dsgnprof", "aaapprof", or "noidear" and you'll see different
lists of courses you can fake-add students to.

Banner Import
---

You can import actual Banner export files with the CSV importer.  You'll
need to change your `BANNER_CSV_PATH` variable and then run the importer:

```bash
SCI_BANNER_CSV_PATH="/path/to/dev/seed/data" ./bin/sci import-csv
```

Once you have populated the database, you can fake a login as any real users to
see what courses are available for attaching students.
