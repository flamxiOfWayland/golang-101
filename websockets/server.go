package main

import (
	"flag"
	"log/slog"
	"net/http"

	"github.com/flamxiOfWayland/golang-101/websockets/server"
)

func main() {
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()

	slog.Info("server is on")
	http.HandleFunc("/ws", server.ServeWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		slog.Error("stopped serving", "err", err)
	}
}
