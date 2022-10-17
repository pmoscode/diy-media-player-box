package schema

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	CardId string
}
