package schema

import (
	"gorm.io/gorm"
	"time"
)

type AudioBook struct {
	gorm.Model
	Title       string
	LastPlayed  time.Time
	CardId      string
	TimesPlayed int
	TrackList   []*AudioTrack
	IsMusic     bool `gorm:"default:false"`
}
