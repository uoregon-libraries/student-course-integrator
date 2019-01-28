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

- Point out places in sci.conf that *must* be overridden
- Seed data for dev so users can test out functionality

### Set up the database

You'll need to set up a database and user for SCI to store its faculty/course
association.  For development, this is trivial by using the supplied
docker-compose configuration:

```bash
docker-compose up -d
```

That would generate a database with username, password, and database name of
"sci" (which match the defaults in the example configuration).

Use `make dbmigrate` to run the goose migrations.

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

TODO:
---

- Implement `make dbmigrate`!!
- Continue fixing this document - still need to reintegrate the bits below
  about debug URLs, running the server, using entr and the devloop script, ...
- Maybe it's time to get docker-compose wrapping this project, too - could ease
  configuration as well as the dev loop script.
- Build some dummy seed data!

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
