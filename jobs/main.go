package main

import (
	sdk "github.com/go-gilbert/gilbert-sdk"
	"github.com/gorilla/websocket"
)

type params struct {
	Address string
	Timeout sdk.Period
}

func NewReloadServerAction(scope sdk.ScopeAccessor, ap sdk.ActionParams) (sdk.ActionHandler, error) {
	p := params{
		Address: "localhost:4022",
		Timeout: sdk.Period(1000),
	}

	if err := ap.Unmarshal(&p); err != nil {
		return nil, err
	}

	return &ReloadServerAction{
		params: p,
		msgs:   make(chan bool, 5),
		upg: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}, nil
}

func GetPluginName() string {
	return "live-reload"
}

func GetPluginActions() sdk.Actions {
	return sdk.Actions{
		"start-server": NewReloadServerAction,
	}
}

func main() {}
