package model

import "fmt"

type Album struct {
	Name  string
	Songs []Song
}

func (a Album) String() string {
	s := fmt.Sprintf("Name: %s\n", a.Name)
	for _, song := range a.Songs {
		s = s + fmt.Sprintf("Song: %s\n", song)
	}

	return s
}
