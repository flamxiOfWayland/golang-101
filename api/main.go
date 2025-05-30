package main

import (
	"fmt"
	"log/slog"
)

func main() {
	controller := NewController(
		WithAddressLogs(),
		WithAllowedClients(
			[]string{
				"0.0.0.1:80",
				// "0.0.0.0:80",
			},
		),
	)

	fmt.Println("API Gateway server started")
	if err := controller.Run(); err != nil {
		slog.Error("Error while serving", "err", err)
	}
}
