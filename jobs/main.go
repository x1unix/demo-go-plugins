package main

import (
	sdk "github.com/go-gilbert/gilbert-sdk"
)

const defaultAddr = "localhost:4800"

type params struct {
	Address string
	Timeout sdk.Period
}

// GetPluginName returns plugin name
func GetPluginName() string {
	return "live-reload"
}

// GetPluginActions returns plugin action handlers list
func GetPluginActions() sdk.Actions {
	return sdk.Actions{
		"start-server": NewReloadServerAction,
		"trigger":      NewReloadTriggerAction,
	}
}

func main() {}
