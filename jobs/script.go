package main

import (
	"fmt"

	sdk "github.com/go-gilbert/gilbert-sdk"
)

const js = `
(function(w) {
	const ADDR = "ws://%s/connect";
	const TIMEOUT = %d;

	console.log("live-reload: connecting to " + ADDR + " ...");
	const socket = new WebSocket(ADDR);
	w.addEventListener('beforeunload', () => socket.close());
	socket.onopen = () => {
		console.log("live-reload: successfully connected to " + ADDR);
	};

	socket.onerror = (e) => console.error("live-reload: error", e);
	socket.onclose = () => console.log("live-reload: disconnected from server"); 

	socket.onmessage = (event) => {
		try {
			const msg = JSON.parse(event.data);
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
		} catch (err) {
			console.error("live-reload: failed to parse message, " + err.message);
		}
	};
})(window)
`

func getConnectionScript(addr string, timeout sdk.Period) []byte {
	out := fmt.Sprintf(js, addr, timeout)
	return []byte(out)
}
