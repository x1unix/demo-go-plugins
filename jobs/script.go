package main

import (
	"fmt"

	sdk "github.com/go-gilbert/gilbert-sdk"
)

const js = `
(function(w) {
	const ADDR = "http://%s/listen";
	const TIMEOUT = %d;

	console.log("live-reload: listening for changes from " + ADDR + " ...");
	w.fetch(ADDR)
	.then(r => r.json())
	.then(msg => {
		switch (msg.type) {
		case "reload":
			console.info("live-reload: reloading...");
			setTimeout(() => w.location.reload(), TIMEOUT);
			break;
		case "shutdown":
			console.info("live-reload: server sent shutdown event");
			socket.close();
			break;
		default:
			console.warn("live-reload: unknown message type", msg);
			break;
		}
	}).catch(err => console.error("live-reload: disconnected from server (" + err.message + ")"))
})(window)
`

func getConnectionScript(addr string, timeout sdk.Period) []byte {
	out := fmt.Sprintf(js, addr, timeout)
	return []byte(out)
}
