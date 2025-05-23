package rapper

import (
	"github.com/flamxiOfWayland/golang-101/model"
)

type Buta struct {
	Name     string `json:"foo"`
	From     string `json:"bar"`
	Albums   []model.Album
	Label    model.Label
	BestSong model.Song
}

func (b Buta) GetAlbums() []model.Album {
	return b.Albums
}

func (b Buta) GetLabel() model.Label {
	return b.Label
}

func (b Buta) GetBestSong() model.Song {
	return b.BestSong
}

func (b *Buta) ChangeBestSong(s model.Song) {
	b.BestSong = s
}

type optFunc func(b *Buta) error

func CreateButa(opts ...optFunc) (*Buta, error) {
	buta := &Buta{
		Name: "Betim",
		From: "Ferizaj",
	}
	for _, opt := range opts {
		if err := opt(buta); err != nil {
			return nil, err
		}
	}
	return buta, nil
}

func WithLabel(l model.Label) optFunc {
	return func(b *Buta) error {
		b.Label = l
		return nil
	}
}

func WithBestSong(s model.Song) optFunc {
	return func(b *Buta) error {
		b.BestSong = s
		return nil
	}
}

func WithAlbum(a model.Album) optFunc {
	return func(b *Buta) error {
		b.Albums = append(b.Albums, a)
		return nil
	}
}
