package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/flamxiOfWayland/golang-101/api/ping"
)

func main() {
	var host string
	var defaultHost string = "0.0.0.0"
	var port string
	var defaultPort string = "80"
	flag.StringVar(&host, "host", defaultHost, fmt.Sprintf("host name, default: %s", defaultHost))
	flag.StringVar(&port, "port", defaultPort, fmt.Sprintf("port name, default: %s", defaultPort))

	var method string
	var defaultMethod string = "get"
	flag.StringVar(&method, "method", defaultMethod, fmt.Sprintf("method name, default: %s", defaultMethod))

	flag.Parse()

	url := fmt.Sprintf("http://%s:%s/ping", host, port)
	if method == "get" {
		resp, err := http.Get(url)
		if err != nil {
			slog.Error("invalid request sent", "err", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			slog.Error("non ok response", "status", resp.Status)
			return
		}

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("invalid response", "err", err)
			return
		}
		defer resp.Body.Close()

		slog.Info("sucessful request", "data", raw)
	}

	if method == "post" {
		def := ping.Default{
			Who:  "goClient",
			What: "testing",
		}
		raw, err := json.Marshal(def)
		if err != nil {
			panic(err)
		}
		data := bytes.NewBuffer(raw)
		resp, err := http.Post(url, "application/json", data)
		if err != nil {
			slog.Error("invalid request sent", "err", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			slog.Error("non ok response", "status", resp.Status)
			return
		}

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("invalid response", "err", err)
			return
		}
		defer resp.Body.Close()
		slog.Info("sucessful request", "data", responseData)
	}
}
