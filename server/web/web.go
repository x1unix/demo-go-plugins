package web

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"github.com/x1unix/demo-go-plugins/server/config"
	"net/http"
)

var server *http.Server

const staticDir = "public"

// Load loads and starts http server
func Load() error {
	server = createServer()
	logrus.Infof("HTTP server is listening '%s'", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %s", err)
	}
	return nil
}

func createServer() *http.Server {
	r := mux.NewRouter()
	registerHandlers(r)
	mw := negroni.New()
	mw.Use(negroni.NewRecovery())
	mw.Use(negroni.NewStatic(http.Dir(staticDir)))
	mw.UseHandler(r)
	return &http.Server{Addr: config.Main.Listen}
}

// Shutdown shuts down http server
func Shutdown() error {
	if server != nil {
		return server.Shutdown(context.Background())
	}

	server = nil
	return nil
}
