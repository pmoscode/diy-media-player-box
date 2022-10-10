package schema

import "gorm.io/gorm"

type AudioTrack struct {
	gorm.Model
	AudioBookID uint
	Track       string
	Title       string
	Length      int
	FileName    string
}
