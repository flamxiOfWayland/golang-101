package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/flamxiOfWayland/golang-101/chat/common"
	"github.com/gorilla/websocket"
)

type Controller struct {
	URL          url.URL
	User         string
	Room         string
	Ctx          context.Context
	ConnNotifier chan<- bool
	From         chan<- common.Message
	To           <-chan common.Message
	Cmd          <-chan common.CommandRequest
	Conn         *websocket.Conn
	Status       *common.Status
}

func (c *Controller) Connect() error {
	header := http.Header{
		"User": []string{c.User},
	}
	conn, resp, err := websocket.DefaultDialer.Dial(c.URL.String(), header)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("cannot connect to server, user: %s - already connected", c.User)
	}
	c.Conn = conn
	c.Status = &common.Status{
		State: common.NormalState,
	}
	return nil
}

func (c *Controller) Listen() error {
	defer c.Conn.Close()
	defer func() error {
		if err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done, bye")); err != nil {
			return err
		}
		slog.Debug("controller", "shutdown", "closing connection")
		return nil
	}()
	go c.HandleInput()
	go c.HandleOutput()

	<-c.Ctx.Done()
	slog.Debug("controller", "shutdown", "safe")

	return nil
}

func (c *Controller) HandleInput() {
	for {
		select {
		case <-c.Ctx.Done():
			//TODO: while default case is blocking this case won't be hit
			slog.Debug("controller", "handle input shutdown", "safe")
			return
		default:
			msgType, msgRaw, err := c.Conn.ReadMessage()
			if err != nil {
				c.ConnNotifier <- true
				slog.Debug("controller", "handle input shutdown", "safe")
				return
			}
			if msgType != websocket.BinaryMessage {
				return
			}

			go c.HandleInputMessage(msgRaw)
		}
	}
}

func (c *Controller) HandleInputMessage(raw []byte) {
	var msg common.Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}
	// TODO: when status is being sent, its marshalled as empty message
	if msg.User != "" || msg.Data != "" {
		c.From <- msg
		return
	}

	if err := json.Unmarshal(raw, c.Status); err != nil {
		slog.Debug("client", "data input unmarshalling", "failed")
		return
	}
}

func (c *Controller) HandleOutput() {
	for {
		select {
		case <-c.Ctx.Done():
			slog.Debug("controller", "handle output shutdown", "safe")
			return
		case data := <-c.Cmd:
			data.User = c.User
			writeMessage(c.Conn, data)
		case data := <-c.To:
			data.User = c.User
			writeMessage(c.Conn, data)
		}
	}
}

func writeMessage[T common.Message | common.CommandRequest](conn *websocket.Conn, data T) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, raw); err != nil {
		return err
	}
	return nil
}

type OptFunc func(*Controller) error

func WithContext(ctx context.Context) OptFunc {
	return func(c *Controller) error {
		c.Ctx = ctx
		return nil
	}
}
func WithConnectionNotifier(ch chan<- bool) OptFunc {
	return func(c *Controller) error {
		c.ConnNotifier = ch
		return nil
	}
}

func WithReceivingMessageChannel(ch chan<- common.Message) OptFunc {
	return func(c *Controller) error {
		c.From = ch
		return nil
	}
}

func WithSendingMessageChannel(ch <-chan common.Message) OptFunc {
	return func(c *Controller) error {
		c.To = ch
		return nil
	}
}

func WithSendingCommandsChannel(ch <-chan common.CommandRequest) OptFunc {
	return func(c *Controller) error {
		c.Cmd = ch
		return nil
	}
}

func WithUser(user string) OptFunc {
	return func(c *Controller) error {
		c.User = user
		return nil
	}
}

func WithRoom(room string) OptFunc {
	return func(c *Controller) error {
		if room == "" {
			return errors.New("invalid room value")
		}
		c.Room = room
		return nil
	}
}

func NewController(host string, opts ...OptFunc) (*Controller, error) {
	if host == "" {
		return nil, errors.New("invalid host value")
	}
	c := &Controller{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	c.URL = url.URL{
		Scheme: "ws", Host: host, Path: c.Room,
	}
	return c, nil
}
