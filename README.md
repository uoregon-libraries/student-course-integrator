Setup
--

- (Install Go)[https://golang.org/dl/]
- Grab vgo: `go get -u golang.org/x/vgo`
- `make`

If you've installed [entr](http://www.entrproject.org/), you can speed up your
development loop by running [`devloop.sh`](./devloop.sh), which runs
[`makerun.sh`](./makerun.sh) whenever `entr` detects a change to any file or
directory under `src/`.
