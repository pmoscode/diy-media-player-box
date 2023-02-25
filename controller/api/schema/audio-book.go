package schema

import (
	"time"
)

type AudioBookFull struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	LastPlayed  time.Time     `json:"lastPlayed"`
	Card        *Card         `json:"card"`
	TimesPlayed int           `json:"timesPlayed"`
	TrackList   []*AudioTrack `json:"trackList"`
	IsMusic     bool          `json:"isMusic"`
}

type AudioBookUi struct {
	Title   string `json:"title"`
	Card    *Card  `json:"card"`
	IsMusic bool   `json:"isMusic"`
}
