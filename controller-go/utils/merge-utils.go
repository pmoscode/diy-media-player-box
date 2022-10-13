package utils

import "controller/database/schema"

func MergeAudioBook(dest *schema.AudioBook, src schema.AudioBook) error {
	dest.CardId = src.CardId
	dest.Title = src.Title

	return nil
}
