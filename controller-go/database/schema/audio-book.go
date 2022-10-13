package schema

import (
	"gorm.io/gorm"
	"time"
)

type AudioBook struct {
	gorm.Model
	Title       string
	LastPlayed  time.Time
	CardId      *Card
	TimesPlayed int
	TrackList   []*AudioTrack
}
