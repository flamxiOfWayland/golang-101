package client

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/flamxiOfWayland/golang-101/chat/common"
	"github.com/jroimartin/gocui"
)

type UI struct {
	user       string
	room       string
	gui        *gocui.Gui
	prompt     *gocui.View
	to         chan<- common.Message
	chat       *gocui.View
	from       <-chan common.Message
	header     *gocui.View
	ctx        context.Context
	connection <-chan bool
	cmd        chan<- common.CommandRequest
	status     *common.Status
}

func (ui *UI) Start(status *common.Status) error {
	ui.status = status
	go ui.ChatMaintainer()
	if err := ui.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	ui.gui.Close()

	slog.Debug("ui", "shutdown", "safe")
	return nil
}

func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *UI) setLayout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	if v, err := gui.SetView("header", -1, -1, maxX, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Frame = true
		v.Autoscroll = false
		v.Editable = true
		header := fmt.Sprintf("hey %s, feel free to express yourself!", ui.user)
		v.Write([]byte(header))
		ui.header = v
	}

	if v, err := gui.SetView("chat", -1, 1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Frame = true
		v.Title = "chat"
		v.Autoscroll = true
		v.Editable = true
		ui.chat = v
	}

	if v, err := gui.SetView("cmd", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = true
		promptEditor := ui.PromptEditor()
		v.Editor = gocui.EditorFunc(promptEditor)
		v.Editable = true
		v.Clear()
		v.Title = "prompt"
		ui.prompt = v
	}

	gui.SetViewOnTop("header")
	gui.SetViewOnTop("chat")
	gui.SetViewOnTop("cmd")
	if _, err := gui.SetCurrentView("cmd"); err != nil {
		return err
	}
	return nil
}

// TODO: need context
func CreateUI(user, room string, ctx context.Context, connection <-chan bool, from <-chan common.Message, to chan<- common.Message, cmd chan<- common.CommandRequest) (*UI, error) {
	var err error
	ui := &UI{
		user:       user,
		room:       room,
		to:         to,
		from:       from,
		ctx:        ctx,
		connection: connection,
		cmd:        cmd,
	}
	if ui.gui, err = gocui.NewGui(gocui.OutputNormal); err != nil {
		return nil, err
	}

	ui.gui.SetManagerFunc(ui.setLayout)
	if err := ui.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return nil, err
	}

	return ui, nil
}

func (ui *UI) update() {
	ui.gui.Update(func(g *gocui.Gui) error { return nil })
}

type EditorFunc func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier)

func (ui *UI) ChatMaintainer() {
	for {
		select {
		case <-ui.ctx.Done():
			slog.Debug("ui chat maintainer", "shutdown", "safe")
			return
		case <-ui.connection:
			ui.gui.Update(func(g *gocui.Gui) error {
				return gocui.ErrQuit
			})
			slog.Debug("ui", "connection", "shutdown")
			return
		case message := <-ui.from:
			ui.chat.Write([]byte(ui.messageRenderer(message, false, false)))
			ui.update()
		}
	}
}

func (ui *UI) PromptEditor() EditorFunc {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case ch != 0 && mod == 0:
			v.EditWrite(ch)
		case key == gocui.KeySpace:
			v.EditWrite(' ')
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
			v.EditDelete(true)
		case key == gocui.KeyDelete:
			v.EditDelete(false)
		case key == gocui.KeyEnter:
			raw := v.Buffer()
			if len(raw) == 0 {
				return
			}
			v.Clear()
			v.EditNewLine()
			data := strings.Trim(raw, "\n")
			go ui.handlePrompt(data)
		}
	}
}

func (ui *UI) handlePrompt(raw string) {
	if strings.Contains(raw, "cmd:") {
		if ui.status.State == common.MutedState {
			return
		}
		// a command is expected to be: "cmd:COMMNAD:USER"
		rawCommand := strings.Split(raw, ":")
		if len(rawCommand) != 3 {
			return
		}
		go ui.handlePromptCommand(rawCommand[1], rawCommand[2])
		return
	}
	go ui.handlePromptMessage(raw)
}

func (ui *UI) handlePromptCommand(cmd, target string) {
	actualCmd := common.CommandFromString(cmd)
	if actualCmd == common.InvalidCmd {
		return
	}
	cmdReq := common.CommandRequest{
		User: ui.user,
		Command: common.Command{
			Target:  target,
			Command: actualCmd,
		},
	}
	ui.cmd <- cmdReq
}

func (ui *UI) handlePromptMessage(raw string) {
	data := common.Message{
		Data: raw,
		User: ui.user,
	}
	final := ui.messageRenderer(data, true, true)
	if ui.status.State == common.MutedState {
		final = fmt.Sprintf("[on mute]: %s", final)
	}
	ui.to <- data
	ui.chat.Write([]byte(final))
	ui.update()
}

func (ui *UI) messageRenderer(m common.Message, self, sameLine bool) string {
	if !self {
		data := fmt.Sprintf("%s:\t%s\n", m.User, m.Data)
		return data
	}
	if !sameLine {
		maxX, _ := ui.gui.Size()
		var x int = 0
		if len(m.User)+1 > len(m.Data) {
			x = len(m.User) - 1
		}
		x = maxX - len(m.Data) - 2 - x
		user := strings.Join([]string{strings.Repeat(" ", x), m.User, ":"}, "")
		message := strings.Join([]string{"\n\t", strings.Repeat(" ", x), m.Data}, "")
		data := strings.Join([]string{user, message, "\n"}, "")
		return data
	} else {
		data := fmt.Sprintf("* %s:\t%s\n", m.User, m.Data)
		return data
	}
}
