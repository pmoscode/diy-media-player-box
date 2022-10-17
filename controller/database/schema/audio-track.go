package schema

import "gorm.io/gorm"

type AudioTrack struct {
	gorm.Model
	AudioBookID uint
	Track       uint
	Title       string
	Length      string
	FileName    string
}
