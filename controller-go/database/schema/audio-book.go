package schema

import (
	"gorm.io/gorm"
)

type AudioBook struct {
	gorm.Model
	Title       string
	LastPlayed  string
	CardId      string
	TimesPlayed int
	TrackList   []AudioTrack
}
