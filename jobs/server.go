package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	sdk "github.com/go-gilbert/gilbert-sdk"
	"github.com/gorilla/websocket"
)

const reloadEndpoint = "/reload"

type event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (e *event) send(conn *websocket.Conn) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

// ReloadServerAction implements action that starts live-reload server
type ReloadServerAction struct {
	upg    websocket.Upgrader
	srv    *http.Server
	msgs   chan event
	params params
}

// Call implements sdk.ActionHandler
func (a *ReloadServerAction) Call(ctx sdk.JobContextAccessor, r sdk.JobRunner) error {
	log := ctx.Log()
	mux := http.NewServeMux()

	script := getConnectionScript(a.params.Address, a.params.Timeout)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", `application/javascript; charset="utf-8"`)
		w.Write(script)
	})

	// websocket connect
	mux.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := a.upg.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("live-reload: failed to create socket, %s", err)
			a.Cancel(ctx)
			go ctx.Result(err) // Pass error to task runner
			return
		}

		log.Infof("live-reload: connected to %s", conn.RemoteAddr())
		for {
			select {
			case e := <-a.msgs:
				log.Debug("live.reload: received reload signal")
				if err := e.send(conn); err != nil {
					log.Errorf("live-reload: failed to send reload signal, %s", err)
				}
			case <-ctx.Context().Done():
				log.Debug("live.reload: received shutdown signal")
				e := event{Type: "shutdown"}
				if err := e.send(conn); err != nil {
					log.Errorf("live-reload: failed to send stop signal to consumer, %s", err)
				}
			}
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
		close(a.msgs)
		return fmt.Errorf("failed to start reload server, %s", err)
	}

	return nil
}

// Cancel shutdowns server
func (a *ReloadServerAction) Cancel(ctx sdk.JobContextAccessor) error {
	ctx.Log().Debug("live-reload: stop server")
	close(a.msgs)
	return a.srv.Shutdown(ctx.Context())
}

// NewReloadServerAction creates a new live-reload start handler
func NewReloadServerAction(scope sdk.ScopeAccessor, ap sdk.ActionParams) (sdk.ActionHandler, error) {
	p := params{
		Address: defaultAddr,
		Timeout: sdk.Period(1000),
	}

	if err := ap.Unmarshal(&p); err != nil {
		return nil, err
	}

	return &ReloadServerAction{
		params: p,
		msgs:   make(chan event),
		upg: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}, nil
}
