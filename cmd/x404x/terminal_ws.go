package main

// terminal_ws.go — WebSocket terminal endpoint for the dashboard.
//
// handleTerminalWS exposes a minimal read-eval loop over a WebSocket so the
// dashboard's xterm.js widget has a live endpoint to talk to. It was missing
// during the audit (dashboard_run.go referenced it) and is recreated here.

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ruby570bocadito/x404x/internal/appstate"
)

var termUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handleTerminalWS returns an http.HandlerFunc that upgrades to a WebSocket and
// serves a small interactive terminal (echo + status commands).
func handleTerminalWS(state *appstate.AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := termUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		_ = conn.WriteMessage(websocket.TextMessage, []byte("X404X web terminal connected.\n"))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("Commands: help, status, exit\n"))

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var out string
			switch string(msg) {
			case "exit", "quit":
				_ = conn.WriteMessage(websocket.TextMessage, []byte("bye\n"))
				return
			case "help":
				out = "Commands: help, status, exit\n"
			case "status":
				if state != nil && state.Bridge != nil {
					out = fmt.Sprintf("bridge connected: %v\n", state.Bridge.Connected())
				} else {
					out = "no state\n"
				}
			default:
				out = fmt.Sprintf("unknown command: %q (try 'help')\n", string(msg))
			}

			if err := conn.WriteMessage(mt, []byte(out)); err != nil {
				return
			}
		}
	}
}
