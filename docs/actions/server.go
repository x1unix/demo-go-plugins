package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	sdk "github.com/go-gilbert/gilbert-sdk"
)

const reloadEndpoint = "/reload"

type event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// ReloadServerAction implements action that starts live-reload server
type ReloadServerAction struct {
	srv    *http.Server
	msgs   chan event
	params params
}

// Call implements sdk.ActionHandler
func (a *ReloadServerAction) Call(ctx sdk.JobContextAccessor, r sdk.JobRunner) (err error) {
	log := ctx.Log()
	mux := http.NewServeMux()

	script := getConnectionScript(a.params.Address, a.params.Timeout)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", `application/javascript; charset="utf-8"`)
		w.WriteHeader(http.StatusOK)
		w.Write(script)
	})

	// change listener
	mux.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Log().Debugf("live-reload: connected to %s", r.RemoteAddr)
		for {
			e, ok := <-a.msgs
			if !ok {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			ctx.Log().Debugf("live.reload: sending refresh to %s", r.RemoteAddr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(e); err != nil {
				ctx.Log().Errorf("live-reload: failed to send response (%s)", err)
			}
			return
		}
	})

	// reload notifier
	mux.HandleFunc(reloadEndpoint, func(w http.ResponseWriter, r *http.Request) {
		log.Debug("live-reload: send reload signal")
		a.msgs <- event{Type: "reload"}
		w.WriteHeader(http.StatusNoContent)
	})

	a.srv = &http.Server{
		Addr:    a.params.Address,
		Handler: mux,
	}

	log.Infof("live-reload: started server at '%s'", a.params.Address)
	if err := a.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start reload server, %s", err)
	}

	return nil
}

// Cancel shutdowns server
func (a *ReloadServerAction) Cancel(ctx sdk.JobContextAccessor) error {
	ctx.Log().Debug("live-reload: stop server")
	return a.srv.Shutdown(ctx.Context())
}

// NewReloadServerAction creates a new live-reload start handler
func NewReloadServerAction(scope sdk.ScopeAccessor, ap sdk.ActionParams) (sdk.ActionHandler, error) {
	p, err := parseParams(scope, ap)
	if err != nil {
		return nil, err
	}

	return &ReloadServerAction{
		params: p,
		msgs:   make(chan event, 5),
	}, nil
}

func parseParams(scope sdk.ScopeAccessor, ap sdk.ActionParams) (params, error) {
	p := params{
		Address: defaultAddr,
		Timeout: sdk.Period(1000),
	}

	if err := ap.Unmarshal(&p); err != nil {
		return p, err
	}

	if err := scope.Scan(&p.Address); err != nil {
		return p, err
	}

	return p, nil
}
