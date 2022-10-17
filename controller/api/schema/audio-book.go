package schema

import (
	"time"
)

type AudioBook struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	LastPlayed  time.Time     `json:"lastPlayed"`
	Card        *Card         `json:"card"`
	TimesPlayed int           `json:"timesPlayed"`
	TrackList   []*AudioTrack `json:"trackList"`
}
