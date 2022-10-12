package schema

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	AudioBookID uint
	CardId      string
}
