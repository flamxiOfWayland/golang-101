package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/flamxiOfWayland/golang-101/chat/common"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Conn struct {
	peer      chan []byte
	broadcast chan []byte
	meta      common.Conn
	self      *websocket.Conn
	others    *Hub
	status    common.Status
}

func NewConn(w http.ResponseWriter, r *http.Request, h *Hub) (*Conn, error) {
	user := r.Header["User"][0]
	if user == "" {
		user = r.Host
	}
	if h.conns[user] != nil {
		w.WriteHeader(http.StatusConflict)
		return nil, fmt.Errorf("server: user: %s already connected!", user)
	}
	slog.Info("server", "new client", user)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	c := &Conn{
		peer:      make(chan []byte),
		broadcast: make(chan []byte),
		others:    h,
		meta: common.Conn{
			User: user,
			Role: common.Admin,
		},
		status: common.ConstructStatus(common.NormalState),
	}
	c.self = ws

	go c.Read()
	go c.Write()

	return c, nil
}

// TODO: message spawn this goroutine
func (c *Conn) ProcessMessage(data []byte) {
	msg := common.Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		slog.Debug("serever process message", "context", err)
		return
	}
	if msg.Data == "" {
		return
	}
	c.broadcast <- data
}

// TODO: message spawn this goroutine
func (c *Conn) ProcessCommand(data []byte) {
	cmd := common.CommandRequest{}
	if err := json.Unmarshal(data, &cmd); err != nil {
		slog.Debug("serever", "process command", err)
		return
	}

	switch cmd.Command.Command {
	case common.InvalidCmd:
		return
	case common.KickCmd:
		if c.meta.User == cmd.Command.Target {
			slog.Debug("server", "process kick command", "someone dumb wanted to kick himself", "who?", cmd.User)
			return
		}
		if c.meta.Role != common.Admin {
			slog.Debug("server", "process kick command", "but user is not admin")
			return
		}
		target, ok := c.others.conns[cmd.Command.Target]
		if !ok {
			slog.Debug("server", "process kick command", "target user is not connected", "user", cmd.Command.Target)
			return
		}
		target.self.Close()
		delete(c.others.conns, cmd.Command.Target)
		slog.Debug("server", "process kick commnad", "succesfull", "admin", cmd.User, "kicked user", cmd.Command.Target)

	case common.MuteCmd:
		if c.meta.User == cmd.Command.Target {
			slog.Debug("server", "process mute command", "someone dumb wanted to mute himself", "who?", cmd.User)
			return
		}
		if c.meta.Role != common.Admin && c.meta.Role != common.Moderator {
			slog.Debug("server", "process mute command", "but user do not have permissions")
			return
		}
		target, ok := c.others.conns[cmd.Command.Target]
		if !ok {
			slog.Debug("server", "process mute command", "target user is not connected", "user", cmd.Command.Target)
			return
		}
		if target.status.State == common.MutedState {
			slog.Debug("server", "process mute command", "target user is already muted, resetting", "user", cmd.Command.Target)
			target.status.Reset()
			return
		}
		target.status = common.ConstructStatus(common.MutedState)
		go target.StatusUpdate()

		slog.Debug("server", "process mute commnad", "succesfull", "admin", cmd.User, "muted user", cmd.Command.Target)
	}
	return
}

func (c *Conn) Read() {
	defer c.self.Close()
	c.self.SetReadLimit(512)
	c.self.SetReadDeadline(time.Now().Add(pongWait))
	c.self.SetPongHandler(func(string) error { c.self.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msgType, data, err := c.self.ReadMessage()
		if err != nil {
			return
		}
		if msgType != websocket.BinaryMessage {
			continue
		}
		switch c.status.State {
		case common.BannedState:
			fallthrough
		case common.NormalState:
			go c.ProcessMessage(data)
			go c.ProcessCommand(data)
		case common.MutedState:
			go c.ProcessCommand(data)
			// if time.Now().Sub(c.status.Start) >= c.status.Duration {
			// 	slog.Debug("server", "user muted status expired", c.meta.User)
			// 	c.status = common.ConstructStatus(common.NormalState)
			// 	go c.StatusUpdate(c.status)
			// 	go c.ProcessMessage(data)
			// 	go c.ProcessCommand(data)
			// } else {
			// 	slog.Debug("server", "user is muted", c.meta.User)
			// 	go c.StatusUpdate(c.status)
			// }
		}
	}
}

func (c *Conn) StatusUpdate() {
	// send the status to the peer
	data, err := json.Marshal(c.status)
	if err != nil {
		slog.Debug("server", "status update failed", err)
		return
	}
	c.peer <- data

	// handle restart and status update
	c.status.Wait()
	c.status = common.ConstructStatus(common.NormalState)
	data, err = json.Marshal(c.status)
	if err != nil {
		slog.Debug("server", "status update failed", err)
		return
	}
	c.peer <- data
}

func (c *Conn) Write() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		c.self.Close()
	}()
	for {
		select {
		case data := <-c.broadcast:
			for _, client := range c.others.conns {
				if client.meta.User == c.meta.User {
					continue
				}
				client.self.SetWriteDeadline(time.Now().Add(writeWait))
				if err := client.self.WriteMessage(websocket.BinaryMessage, data); err != nil {
					return
				}
				slog.Debug("server broadcasting", "user", c.meta.User, "said", data)
			}
		case data := <-c.peer:
			c.self.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.self.WriteMessage(websocket.BinaryMessage, data); err != nil {
				return
			}
		case <-pingTicker.C:
			c.self.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.self.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	addr  string
	room  string
	conns map[string]*Conn
	ctx   context.Context
}

// TODO: this might not be terminated correctly
func (h *Hub) ListenAndServe() {
	http.HandleFunc(h.room, h.ServeWs)
	if err := http.ListenAndServe(h.addr, nil); err != nil {
		slog.Error("listening terminated", "error", err)
		return
	}
	slog.Debug("listening terminated")
}

func (h *Hub) HandleCloseConnection(conn *Conn) {
	conn.self.SetCloseHandler(func(code int, text string) error {
		slog.Debug("client", conn.meta.User, "left", "info", text)
		delete(h.conns, conn.meta.User)
		return nil
	})
}

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	c, err := NewConn(w, r, h)
	if err != nil {
		slog.Error("new connection", "error", err)
		return
	}
	h.conns[c.meta.User] = c
	h.HandleCloseConnection(c)
}

func NewHub(addr, room string, ctx context.Context) *Hub {
	h := &Hub{
		conns: make(map[string]*Conn),
		ctx:   ctx,
		addr:  addr,
		room:  room,
	}
	return h
}
