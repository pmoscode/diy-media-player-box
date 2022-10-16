package utils

import (
	uiSchema "controller/api/schema"
	dbSchema "controller/database/schema"
)

func ConvertAudioBookUiToDb(src *uiSchema.AudioBook) *dbSchema.AudioBook {
	dest := &dbSchema.AudioBook{}

	if src.Card != nil {
		dest.CardId = src.Card.CardId
	}
	dest.Title = src.Title

	return dest
}

func ConvertAudioBookDbToUi(src *dbSchema.AudioBook) *uiSchema.AudioBook {
	dest := &uiSchema.AudioBook{}

	dest.Id = int(src.ID)
	if src.CardId != "" {
		dest.Card = &uiSchema.Card{
			CardId: src.CardId,
		}
	}
	dest.Title = src.Title
	dest.LastPlayed = src.LastPlayed
	dest.TimesPlayed = src.TimesPlayed

	if src.TrackList != nil {
		dest.TrackList = ConvertAudioBookTracksDbToUi(src.TrackList)
	}

	return dest
}

func MergeAudioBookUiToDb(dest *dbSchema.AudioBook, src *uiSchema.AudioBook) {
	dest.Title = src.Title
	dest.CardId = src.Card.CardId
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
