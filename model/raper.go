package model

type Rapper interface {
	GetAlbums() []Album
	GetBestSong() Song
	GetLabel() Label
}
