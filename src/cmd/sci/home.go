package main

import (
	"net/http"

	"github.com/uoregon-libraries/gopkg/tmpl"
)

// homeHandler encapsulates basic data and functionality for handling input and
// rendering output
type homeHandler struct {
	tmpl *tmpl.Template
}

func hHome() *homeHandler {
	var r = layout.Clone()
	return &homeHandler{
		tmpl: r.MustBuild("home.go.html"),
	}
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.tmpl.Execute(w, nil)
}
