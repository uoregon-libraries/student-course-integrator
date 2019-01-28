Setup (for developers)
--

This is not meant to help end-users.  That's a guide we also need but don't
have time to work on right now.

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

Diehards would recommend using `go get`, but it's sort of magic, it tries to
compile things you don't necessarily want, and it's just a lot clearer if you
do a git clone.

```bash
mkdir -p $GOPATH/src/github.com/uoregon-libraries
git clone git@github.com:uoregon-libraries/student-course-integrator.git \
          $GOPATH/src/github.com/uoregon-libraries/student-course-integrator
cd $GOPATH/src/github.com/uoregon-libraries/student-course-integrator
```

### Get Goose

Get goose for database migrations: `go get -u bitbucket.org/liamstask/goose/...`

Verify it worked and that your paths are all set up: `goose --help`.  You
should get a "usage" blurb.

### Get golint

`golint` is a tool for analyzing your Go code to look for various best practices.

```bash
go get -u golang.org/x/lint/golint
```

### Configure SCI

Copy the example configuration and edit it:

```bash
cp example.conf sci.conf
vim sci.conf
```

You don't have to use vim, of course.  Neovim is also acceptable.

Look over the file carefully, you will need to configure several things to make
SCI run locally.  The example configuration is set up to make development
easier, but there are still details not filled in.  The file is well-commented
to ensure easier understanding of the configuration needs.

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

Run the server
---

### Summary

Read below for details, but this is the quick setup info:

```bash
# Start up the database server if you haven't already done so
docker-compose up -d

# Migrate the database - this must happen at least once when you first start
# the project, then again anytime a migration is added
make dbmigrate

# Populate the database - this must happen **only once**, otherwise any local
# changes to courses or faculty will be destroyed.
docker-compose exec db mysql -usci -psci -Dsci -e "source /tmp/seed.sql"

# Run the SCI http listener
./bin/sci server
```

### Prepare the database

You'll need to set up a database for SCI to store its faculty/course
association.  For development, this is trivial by using docker-compose:
`docker-compose.yml` defines a database with username, password, and database
name of "sci" (which match the defaults in the example configuration).  Simply
run `docker-compose up -d` to start it.

Once your database is up, run `make dbmigrate` to run the goose migrations.

### Verify config

Make sure you've got `sci.conf` set up for development.

### Populate the database

You will need to populate the "sci" database tables with any fake (or real)
courses and user ids.

You can load a small sample of test data by importing all the SQL in `seed.sql`
into your dockerized database container:

```bash
docker-compose exec db mysql -usci -psci -Dsci -e "source /tmp/seed.sql"
```

**Note**: on your first use, or after docker volumes have been destroyed, you'll want to
seed data.  But you don't re-run seeds every time you start up the server.

### Run SCI

Simply type `./bin/sci server`.  The command reads your local `sci.conf`, so it
doesn't need any command-line configuration.

### Log in

Visit `http://localhost:8080` (if your sci.conf kept the default port) and you
should see the app's "not authorized" page, and you'll be logged in as
"dummyuser" in place of a real authorization.

In debug mode, you can fake a login as any other user by visiting the page with
a "debuguser" query argument.  For example, `http://localhost:8080/?debuguser=jechols`.
The page will have a large, visible warning if it is in debug mode to avoid
accidentally pushing debug to production.

If you used the seed data, use the "debuguser" argument to sign in as
"dsgnprof", "aaapprof", or "noidear" and you'll see different lists of courses
you can fake-add students to.

Development loop
---

If you install [entr](http://www.entrproject.org/), you can speed up your
development loop by running [`./scripts/devloop.sh`](./scripts/devloop.sh),
which runs [`makerun.sh`](./scripts/makerun.sh) whenever `entr` detects a
change to any file or directory under `src/`.

Build and test
---

You can build the binary by simply running `make`, but it's advised that you
also validate the code and run tests before pushing anything up:

```bash
make validate
make test
make
```

Banner Import
---

You can also import actual Banner export files with the CSV importer.  You'll
need to change your `BANNER_CSV_PATH` variable and then run the importer:

```bash
SCI_BANNER_CSV_PATH="/path/to/dev/seed/data" ./bin/sci import-csv
```

Once you have populated the database, you can fake a login as any real users to
see what courses are available for attaching students.
