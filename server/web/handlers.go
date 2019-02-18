package web

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/x1unix/demo-go-plugins/server/feed"
	"net/http"
	"strconv"
)

const (
	countParam      = "count"
	afterIDParam    = "after"
	sectionParam    = "section"
	sourceNameParam = "sourceName"
	defaultCount    = 30
)

func newRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Path("/sources").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		OK(response{
			"sources": feed.SourceNames(),
		}, w)
	})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(errors.New("resource not found"), http.StatusNotFound, w)
	})

	r.Path("/sources/{sourceName}/sections").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		src, err := sourceFromRequest(r)
		if err != nil {
			Error(err, http.StatusNotFound, w)
			return
		}

		OK(response{
			"source":   src.Name(),
			"sections": src.Sections(),
		}, w)
	})

	r.Path("/sources/{sourceName}/sections/{section}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		src, err := sourceFromRequest(r)
		if err != nil {
			Error(err, http.StatusNotFound, w)
			return
		}

		section, selector := readSelectorParams(r)
		posts, err := src.GetPosts(section, selector)
		if err != nil {
			Error(err, http.StatusInternalServerError, w)
			return
		}

		OK(response{
			"source":  src.Name(),
			"section": section,
			"posts":   posts,
		}, w)
	})

	return r
}

func sourceFromRequest(r *http.Request) (feed.Source, error) {
	v := mux.Vars(r)
	name := v[sourceNameParam]
	return feed.GetSource(name)
}

func readSelectorParams(r *http.Request) (section string, selector feed.Selector) {
	v := mux.Vars(r)
	section = v[sectionParam]

	q := r.URL.Query()
	if count, err := strconv.ParseUint(q.Get(countParam), 10, 32); err != nil {
		selector.Count = defaultCount
	} else {
		selector.Count = uint(count)
	}

	selector.AfterID = q.Get(afterIDParam)
	return section, selector
}
