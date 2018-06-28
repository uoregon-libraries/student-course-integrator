package sciserver

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/uoregon-libraries/gopkg/logger"
)

// server is a very simple web server encapsulation, mostly so configuration is
// wrapped in a type instead of passed to the listener
type server struct {
	Approot string // Path to the app for getting things like templates and static assets
	Bind    string // Bind address, e.g., ":80" for production
	Debug   bool   // If true, middleware for hacking userid is added
}

func (s *server) Listen() {
	var r = mux.NewRouter()
	if s.Debug {
		r.Use(fakeUserLogin)
	}

	// Static asset server just lets anything in approot/static through to the browser
	var fileServer = http.FileServer(http.Dir(filepath.Join(s.Approot, "static")))
	var fileRouter = r.NewRoute().PathPrefix("/static").Subrouter()
	fileRouter.Use(requestStaticAssetLog)
	fileRouter.NewRoute().Handler(http.StripPrefix("/static", fileServer))

	// Everything else gets middleware to avoid browser caching, and to log the
	// more meaningful requests.  We use a new subrouter to ensure consistency.
	var sub = r.NewRoute().PathPrefix("").Subrouter()

	// If we're in debug mode, we need to hack the user header before anything else
	sub.Use(getUser, nocache, requestLog, mustAuth)
	sub.NewRoute().Path("/").Handler(hHome())
	sub.NewRoute().HandlerFunc(http.NotFound)

	http.Handle("/", r)

	logger.Infof("Listening on %s", s.Bind)
	if err := http.ListenAndServe(s.Bind, nil); err != nil {
		logger.Fatalf("Error starting listener: %s", err)
	}
}
