package common

import "time"

type Role int

const (
	Admin Role = iota
	Moderator
	Guest
)

type Conn struct {
	User string
	Role Role
}

type State int

const (
	NormalState = iota
	MutedState
	BannedState // TODO: this impl. needs a stateful server (a database | redis)
)

type Status struct {
	State State `json:"state"`
	timer time.Timer
}

func (s *Status) Wait() {
	<-s.timer.C
}
func (s *Status) Reset() {
	if s.timer.Stop() {
		s.timer.Reset(time.Duration(10 * time.Second))
	} else {
		s.timer = *time.NewTimer(time.Duration(10 * time.Second))
	}
}

func ConstructStatus(state State) Status {
	switch state {
	case MutedState:
		return Status{
			State: MutedState,
			timer: *time.NewTimer(time.Duration(10 * time.Second)),
		}
	case BannedState:
		return Status{
			State: BannedState,
			timer: *time.NewTimer(time.Duration(30 * time.Second)),
		}
	case NormalState:
		return Status{
			State: NormalState,
		}
	}
	return Status{
		State: NormalState,
	}
}
