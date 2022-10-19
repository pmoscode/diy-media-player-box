package audio

import (
	"audio-player/mqtt"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
)

const DefaultSampleRate = 44100

type Audio struct {
	control             *beep.Ctrl
	lastPlayedUid       uint
	sendStatusMessage   func(message string)
	sendPlayDoneMessage func(id uint)
}

func (a *Audio) checkLastPlayedUidChanged(body *TracksSubscriptionMessage) bool {
	if a.lastPlayedUid != body.Id {
		if a.lastPlayedUid > 0 {
			a.sendPlayDoneMessage(a.lastPlayedUid)
		}

		a.lastPlayedUid = body.Id

		return true
	}

	return false
}

func (a *Audio) OnMessageReceivedPlay(message mqtt.Message) {
	body := TracksSubscriptionMessage{}

	message.ToStruct(&body)

	transition := fmt.Sprintf("%d to %d", a.lastPlayedUid, body.Id)
	uidChanged := a.checkLastPlayedUidChanged(&body)
	a.sendStatusMessage("uid changed: " + transition)

	if uidChanged {
		speaker.Clear()

		var samples []beep.Streamer

		for _, trackPath := range body.TrackList {
			f, err := os.Open(trackPath)
			if err != nil {
				a.sendStatusMessage("Could not open '" + trackPath + "'... DYING!!!")
				log.Fatal(err)
			}

			streamer, format, err := mp3.Decode(f)
			if err != nil {
				a.sendStatusMessage("Could not decode '" + trackPath + "'... DYING!!!")
				log.Fatal(err)
			}

			var stream beep.Streamer

			if DefaultSampleRate != format.SampleRate {
				const sampleRate = beep.SampleRate(DefaultSampleRate)
				stream = beep.Resample(1, format.SampleRate, sampleRate, streamer)
				a.sendStatusMessage("Need to resample: '" + trackPath + "'...")
			} else {
				stream = streamer
				a.sendStatusMessage("No need to resample: '" + trackPath + "'...")
			}

			samples = append(samples, stream)
		}

		samples = append(samples, beep.Callback(func() {
			a.lastPlayedUid = 0
			a.sendStatusMessage("stopped")
			a.sendPlayDoneMessage(body.Id)
		}))

		if len(samples) > 0 {
			sequence := beep.Seq(samples...)
			a.control = &beep.Ctrl{Streamer: sequence, Paused: false}

			speaker.Play(a.control)

			a.sendStatusMessage("playing")
		} else {
			a.sendStatusMessage("no tracks")
		}
	} else {
		if a.control.Paused {
			a.OnMessageReceivedSwitch(message)
		} else {
			a.sendStatusMessage("untouched")
		}
	}
}

func (a *Audio) OnMessageReceivedSwitch(message mqtt.Message) {
	if a.control != nil {
		speaker.Lock()
		a.control.Paused = !a.control.Paused
		speaker.Unlock()

		status := "paused"
		if !a.control.Paused {
			status = "continuing"
		}

		a.sendStatusMessage(status)
	} else {
		a.sendStatusMessage("no audio stream")
	}
}

func (a *Audio) OnMessageReceivedStop(message mqtt.Message) {
	speaker.Clear()

	a.lastPlayedUid = 0
	a.sendStatusMessage("stopped")
}

func NewAudio(statusMessage func(statusMessage string), playDoneMessage func(id uint)) *Audio {
	return &Audio{
		sendStatusMessage:   statusMessage,
		sendPlayDoneMessage: playDoneMessage,
	}
}
