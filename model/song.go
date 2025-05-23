package model

import "fmt"

type Song struct {
	Title string
	Views int
}

func (s Song) String() string {
	return fmt.Sprintf("Title: %s, Views: %d", s.Title, s.Views)
}
