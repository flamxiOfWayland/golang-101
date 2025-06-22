package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"github.com/flamxiOfWayland/golang-101/chat/server"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":8080", "service address")
	room := flag.String("room", "/chat", "chat room")
	verbose := flag.Bool("verbose", false, "verbose logs")
	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, done := context.WithCancel(context.Background())

	hub := server.NewHub(*addr, *room, ctx)
	slog.Info("listening for connections...")

	go hub.ListenAndServe()

	<-interrupt
	done()

	slog.Info("server is shut down")
}
