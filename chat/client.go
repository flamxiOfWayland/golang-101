package main

import (
	"context"
	"flag"
	"log/slog"
	"time"

	"github.com/flamxiOfWayland/golang-101/chat/client"
	"github.com/flamxiOfWayland/golang-101/chat/common"
)

func main() {
	host := flag.String(
		"host",
		":8080",
		"http service address",
	)
	user := flag.String(
		"user",
		"default",
		"user you want to identify with",
	)
	verbose := flag.Bool(
		"verbose",
		false,
		"verbosity level, default info",
	)
	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	ctx, cancel := context.WithCancel(context.Background())

	from := make(chan common.Message)
	to := make(chan common.Message)
	cmd := make(chan common.CommandRequest)
	connection := make(chan bool)

	controller, err := client.NewController(
		*host,
		client.WithContext(ctx),
		client.WithConnectionNotifier(connection),
		client.WithReceivingMessageChannel(from),
		client.WithSendingMessageChannel(to),
		client.WithSendingCommandsChannel(cmd),
		client.WithUser(*user),
		client.WithRoom("/chat"),
	)
	if err != nil {
		slog.Error("creating ws controller", "error", err)
		return
	}
	if err := controller.Connect(); err != nil {
		slog.Error("ws controller", "connecting error", err)
		return
	}

	ui, err := client.CreateUI(*user, controller.URL.String(), ctx, connection, from, to, cmd)
	if err != nil {
		slog.Error("creating ui", "error", err)
		return
	}

	go controller.Listen()
	if err := ui.Start(controller.Status); err != nil {
		slog.Error("ui", "error", err)
	}
	cancel()

	time.Sleep(1 * time.Millisecond)
	slog.Debug("app", "shutdown", "safe")
}
