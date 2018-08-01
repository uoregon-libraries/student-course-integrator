package version

// version is displayed on web pages so we can be sure dev/staging/prod are the
// version we expect them to be
const version = "1.0.0"

func Version() string {
	return version
}
