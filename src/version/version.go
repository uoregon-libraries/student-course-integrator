//go:generate go run gen.go

package version

// version is displayed on web pages so we can be sure dev/staging/prod are the
// version we expect them to be
const version = "1.3.2"

// commit is empty if the commit.go file doesn't get generated - this isn't
// super helpful, but it does allow one to download and build without `make` if
// necessary.
var commit = "xyzzy"

// Version returns a combination of main version string and commit
func Version() string {
	return version + "-" + commit
}
