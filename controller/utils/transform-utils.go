package utils

import (
	uiSchema "controller/api/schema"
	dbSchema "controller/database/schema"
)

func ConvertAudioBookUiToDb(src *uiSchema.AudioBookUi) *dbSchema.AudioBook {
	dest := &dbSchema.AudioBook{}

	if src.Card != nil {
		dest.CardId = src.Card.CardId
	}
	dest.Title = src.Title
	dest.IsMusic = src.IsMusic

	return dest
}

func ConvertAudioBookDbToUi(src *dbSchema.AudioBook) *uiSchema.AudioBookFull {
	dest := &uiSchema.AudioBookFull{}

	dest.Id = int(src.ID)
	if src.CardId != "" {
		dest.Card = &uiSchema.Card{
			CardId: src.CardId,
		}
	}
	dest.Title = src.Title
	dest.LastPlayed = src.LastPlayed
	dest.TimesPlayed = src.TimesPlayed
	dest.IsMusic = src.IsMusic

	if src.TrackList != nil {
		dest.TrackList = ConvertAudioBookTracksDbToUi(src.TrackList)
	}

	return dest
}

func MergeAudioBookUiToDb(dest *dbSchema.AudioBook, src *uiSchema.AudioBookUi) {
	dest.Title = src.Title
	dest.IsMusic = src.IsMusic

	if src.Card != nil {
		dest.CardId = src.Card.CardId
	}
}

func ConvertAudioBookTracksDbToUi(src []*dbSchema.AudioTrack) []*uiSchema.AudioTrack {
	uiAudioTracks := make([]*uiSchema.AudioTrack, 0)

	for _, audioTrack := range src {
		uiTrack := &uiSchema.AudioTrack{
			Track:    audioTrack.Track,
			Title:    audioTrack.Title,
			Length:   audioTrack.Length,
			FileName: audioTrack.FileName,
		}
		uiAudioTracks = append(uiAudioTracks, uiTrack)
	}

	return uiAudioTracks
}

func ConvertCardDbToUi(src *dbSchema.Card) *uiSchema.Card {
	return &uiSchema.Card{
		Id:     int(src.ID),
		CardId: src.CardId,
	}
}
