package audio

import (
	"audio-player/mqtt"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
)

const DefaultSampleRate = 44100

type Audio struct {
	control             *beep.Ctrl
	volume              *effects.Volume
	lastPlayedUid       uint
	sendStatusMessage   func(messageType mqtt.StatusType, message ...any)
	sendPlayDoneMessage func(id uint)
}

func (a *Audio) checkLastPlayedUidChanged(body *mqtt.TracksSubscriptionMessage) bool {
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
	body := mqtt.TracksSubscriptionMessage{}

	message.ToStruct(&body)

	transition := fmt.Sprintf("%d to %d", a.lastPlayedUid, body.Id)
	uidChanged := a.checkLastPlayedUidChanged(&body)
	a.sendStatusMessage(mqtt.Info, "uid changed: ", transition)

	if uidChanged {
		speaker.Clear()

		var samples []beep.Streamer

		for _, trackPath := range body.TrackList {
			f, err := os.Open(trackPath)
			if err != nil {
				a.sendStatusMessage(mqtt.Error, "Could not open '"+trackPath+"'... DYING!!!")
				log.Fatal(err)
			}

			streamer, format, err := mp3.Decode(f)
			if err != nil {
				a.sendStatusMessage(mqtt.Error, "Could not decode '"+trackPath+"'... DYING!!!")
				log.Fatal(err)
			}

			var stream beep.Streamer

			if DefaultSampleRate != format.SampleRate {
				const sampleRate = beep.SampleRate(DefaultSampleRate)
				stream = beep.Resample(1, format.SampleRate, sampleRate, streamer)
				a.sendStatusMessage(mqtt.Info, "Need to resample: '"+trackPath+"'...")
			} else {
				stream = streamer
				a.sendStatusMessage(mqtt.Info, "No need to resample: '"+trackPath+"'...")
			}

			samples = append(samples, stream)
		}

		samples = append(samples, beep.Callback(func() {
			a.lastPlayedUid = 0
			a.sendStatusMessage(mqtt.Info, "stopped")
			a.sendPlayDoneMessage(body.Id)
		}))

		if len(samples) > 0 {
			sequence := beep.Seq(samples...)
			a.control = &beep.Ctrl{
				Streamer: sequence,
				Paused:   false,
			}

			a.volume = &effects.Volume{
				Streamer: a.control,
				Base:     2,
				Volume:   0,
				Silent:   false,
			}

			speaker.Play(a.volume)

			a.sendStatusMessage(mqtt.Info, "playing")
		} else {
			a.sendStatusMessage(mqtt.Info, "no tracks")
		}
	} else {
		if a.control.Paused {
			a.OnMessageReceivedSwitch(message)
		} else {
			a.sendStatusMessage(mqtt.Info, "untouched")
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

		a.sendStatusMessage(mqtt.Info, status)
	} else {
		a.sendStatusMessage(mqtt.Info, "no audio stream")
	}
}

func (a *Audio) OnMessageReceivedStop(message mqtt.Message) {
	speaker.Clear()

	a.lastPlayedUid = 0
	a.sendStatusMessage(mqtt.Info, "stopped")
}

func (a *Audio) OnMessageReceivedVolume(message mqtt.Message) {
	volumeMessage := &mqtt.VolumeChangeSubscriptionMessage{}

	message.ToStruct(volumeMessage)

	speaker.Lock()
	a.volume.Volume += volumeMessage.VolumeOffset
	speaker.Unlock()
}

func NewAudio(statusMessage func(messageType mqtt.StatusType, message ...any), playDoneMessage func(id uint)) *Audio {
	return &Audio{
		sendStatusMessage:   statusMessage,
		sendPlayDoneMessage: playDoneMessage,
	}
}
