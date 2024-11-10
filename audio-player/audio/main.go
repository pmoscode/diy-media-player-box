package audio

import (
	mqtt2 "github.com/pmoscode/go-common/mqtt"
)

func NewAudio(statusMessage func(messageType mqtt2.StatusType, message ...any), playDoneMessage func(id uint)) *Audio {
	return &Audio{
		sendStatusMessage:   statusMessage,
		sendPlayDoneMessage: playDoneMessage,
		currentVolume:       0.0,
	}
}
