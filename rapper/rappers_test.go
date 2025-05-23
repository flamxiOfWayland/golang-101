package rapper_test

import (
	"fmt"
	"testing"

	"github.com/flamxiOfWayland/golang-101/model"
	"github.com/flamxiOfWayland/golang-101/rapper"
	"github.com/stretchr/testify/assert"
)

func TestRappers(t *testing.T) {
	rappers := rapper.Rappers()
	assert.Equal(t, len(rappers), 3)

	lyricalSon := rappers[0].(rapper.LyricalSon)
	for _, album := range lyricalSon.GetAlbums() {
		fmt.Println(album)
	}
	lyricalSon.GetAlbums()
	fmt.Println(lyricalSon.Name)

	buta := rappers[2].(*rapper.Buta)
	fmt.Println(buta.GetBestSong().Title)
	buta.ChangeBestSong(model.Song{
		Title: "telebingo"})
	fmt.Println(buta.GetBestSong().Title)

	assert.True(t, false)
}
