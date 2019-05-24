package main

import (
	"fmt"
	"net/http"

	sdk "github.com/go-gilbert/gilbert-sdk"
)

// ReloadTriggerAction triggers live-reload server
type ReloadTriggerAction struct {
	params params
}

// Call implements sdk.ActionHandler
func (a *ReloadTriggerAction) Call(ctx sdk.JobContextAccessor, r sdk.JobRunner) error {
	tr := &http.Transport{}
	req, err := http.NewRequest(http.MethodPost, "http://"+a.params.Address+reloadEndpoint, nil)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Context().Done()
		tr.CancelRequest(req)
	}()

	client := &http.Client{Transport: tr}
	ctx.Log().Debugf("live-reload: sending refresh signal to '%s'", req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to trigger page reload on server (%s)", err)
	}

	if resp.StatusCode > 300 {
		ctx.Log().Warnf("live-reload: reload server responded with bad status %d", resp.StatusCode)
	}

	return nil
}

// Cancel implements sdk.ActionHandler
func (a *ReloadTriggerAction) Cancel(ctx sdk.JobContextAccessor) error {
	// cancelation already implemented in Call()
	return nil
}

// NewReloadTriggerAction creates a new reload trigger action
func NewReloadTriggerAction(scope sdk.ScopeAccessor, ap sdk.ActionParams) (sdk.ActionHandler, error) {
	p, err := parseParams(scope, ap)
	if err != nil {
		return nil, err
	}

	return &ReloadTriggerAction{
		params: p,
	}, nil
}
