Setup
--

- [Install Go](https://golang.org/dl/)
- Grab vgo: `go get -u golang.org/x/vgo`
- `make`

If you've installed [entr](http://www.entrproject.org/), you can speed up your
development loop by running [`devloop.sh`](./devloop.sh), which runs
[`makerun.sh`](./makerun.sh) whenever `entr` detects a change to any file or
directory under `src/`.

Settings file
---

Copy example.conf to sci.conf and modify it as needed.  SCI will look for this
file at `/etc/sci.conf`, then `./sci.conf`.  You can pass the `-c` flag to
specify a custom location as well, e.g., `./bin/sci -c /tmp/dummysettings.conf`.

All settings can be overridden with environment variables prefixed with "SCI_".
In production, use this to avoid storing sensitive values in `sci.conf`:

    export SCI_DB=sciusername:password@tcp(localhost:port)/scidbname
    export SCI_SESSION_SECRET=blah

Run the server
---

`./bin/sci`
