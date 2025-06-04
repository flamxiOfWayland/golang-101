package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/chat"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
		}
	}()

	var reader chan []byte = make(chan []byte)
	go readStdin(reader)

	for {
		select {
		case <-done:
			return
		case data := <-reader:
			err := c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			return
		}
	}
}

func readStdin(pipe chan<- []byte) {
	for {
		var buffer []byte = make([]byte, 512)
		_, err := os.Stdin.Read(buffer)
		if err != nil {
			panic(err)
		}
		data := strings.TrimSuffix(string(buffer), "\n")
		pipe <- []byte(data)
	}
}
