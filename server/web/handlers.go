package web

import (
	"github.com/gorilla/mux"
	"github.com/x1unix/demo-go-plugins/server/feed"
	"net/http"
)

func registerHandlers(r *mux.Router) {
	r.Path("/sources").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		OK(response{
			"sources": feed.SourceNames(),
		}, w)
	})
}
