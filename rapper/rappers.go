package rapper

import (
	"github.com/flamxiOfWayland/golang-101/model"
	"log"
)

type LyricalSon struct {
	Name     string
	From     string
	Albums   []model.Album
	Label    model.Label
	BestSong model.Song
}

func (l LyricalSon) GetAlbums() []model.Album {
	return l.Albums
}

func (l LyricalSon) GetLabel() model.Label {
	return l.Label
}

func (l LyricalSon) GetBestSong() model.Song {
	return l.BestSong
}

func Rappers() []interface{} {
	var rappers []interface{}

	lyricalSon := LyricalSon{
		Name:  "Festim",
		Label: model.PINT,
		From:  "Prishtine",
		Albums: []model.Album{
			{
				Name: "respect",
				Songs: []model.Song{
					{
						Title: "respect",
						Views: 100,
					},
					{
						Title: "albanian",
						Views: 10000,
					},
				},
			},
			{
				Name: "salihi x ferizi",
				Songs: []model.Song{
					{
						Title: "djathi i zvicres",
						Views: 100,
					},
					{
						Title: "oj hane",
						Views: 10000,
					},
				},
			},
		},
	}

	rappers = append(rappers, lyricalSon)

	unikkatil := CreateUnikkatil("viktor", "prishtine", model.TBA, model.Song{Title: "kejt hajvan", Views: 1_000_000_000})

	rappers = append(rappers, unikkatil)

	buta, err := CreateButa(
		WithLabel(model.Label("rifa")),
		WithAlbum(model.Album{Name: "pranver/ver 2020/2021"}),
		WithAlbum(model.Album{Name: "pranver/ver 2021/2022"}),
		WithAlbum(model.Album{Name: "pranver/ver 2022/2023"}),
		WithAlbum(model.Album{Name: "pranver/ver 2023/2024"}),
		WithBestSong(model.Song{Title: "pluto", Views: 1_000_000_000_000}),
	)

	if err != nil {
		log.Fatal(err)
	}
	rappers = append(rappers, buta)

	return rappers
}
