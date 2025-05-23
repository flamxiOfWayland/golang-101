package rapper

import "github.com/flamxiOfWayland/golang-101/model"

type Unikkatil struct {
	Name     string
	From     string
	Albums   []model.Album
	Label    model.Label
	BestSong model.Song
	HasGun   bool
}

func CreateUnikkatil(name, from string, label model.Label, bestSong model.Song) *Unikkatil {
	return &Unikkatil{
		Name:     name,
		From:     from,
		Label:    label,
		BestSong: bestSong,
	}
}

func (u Unikkatil) GetAlbums() []model.Album {
	return u.Albums
}

func (u Unikkatil) GetLabel() model.Label {
	return u.Label
}

func (u Unikkatil) GetBestSong() model.Song {
	return u.BestSong
}
