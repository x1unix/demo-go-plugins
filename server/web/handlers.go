package web

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/x1unix/demo-go-plugins/server/feed"
	"net/http"
)

func newRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Path("/sources").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		OK(response{
			"sources": feed.SourceNames(),
		}, w)
	})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(errors.New("resource not found"), http.StatusNotFound, w)
	})
	return r
	//r.Path("/sources/{sourceName}/{section}")
}
