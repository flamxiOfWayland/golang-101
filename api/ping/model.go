package ping

import "fmt"

type Default struct {
	Who  string `json:"who"`
	What string `json:"what"`
}

func (d Default) String() string {
	return fmt.Sprintf("Who: %s, What: %s", d.Who, d.What)
}
