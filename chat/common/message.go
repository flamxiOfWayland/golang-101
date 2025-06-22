package common

type CommandEnum int

const (
	InvalidCmd CommandEnum = iota
	MuteCmd
	KickCmd
	BanCmd
)

func CommandFromString(raw string) CommandEnum {
	if raw == "kick" {
		return KickCmd
	}
	if raw == "mute" {
		return MuteCmd
	}
	return InvalidCmd
}

type Command struct {
	Command CommandEnum `json:"cmd"`
	Target  string      `json:"target"`
}

type CommandRequest struct {
	User    string  `json:"user"`
	Command Command `json:"command"`
}

type Message struct {
	Data string `json:"data"`
	User string `json:"user"`
}
