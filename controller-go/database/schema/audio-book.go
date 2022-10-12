package schema

import (
	"gorm.io/gorm"
)

type AudioBook struct {
	gorm.Model
	Title       string
	LastPlayed  string
	CardId      *Card
	TimesPlayed int
	TrackList   []AudioTrack
}
