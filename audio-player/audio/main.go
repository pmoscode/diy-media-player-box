package audio

import (
	"audio-player/mqtt"
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

func (a *Audio) OnMessageReceivedPlay(message mqtt.Message) {
	body := mqtt.TracksSubscriptionMessage{}

	message.ToStruct(&body)

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
}

func (a *Audio) OnMessageReceivedPause(message mqtt.Message) {
	if a.control != nil {
		speaker.Lock()
		a.control.Paused = true
		speaker.Unlock()

		a.sendStatusMessage(mqtt.Info, "paused")
	} else {
		a.sendStatusMessage(mqtt.Info, "no audio stream to pause...")
	}
}

func (a *Audio) OnMessageReceivedResume(message mqtt.Message) {
	if a.control != nil {
		speaker.Lock()
		a.control.Paused = false
		speaker.Unlock()

		a.sendStatusMessage(mqtt.Info, "continuing")
	} else {
		a.sendStatusMessage(mqtt.Info, "no audio stream to continue...")
	}
}

func (a *Audio) OnMessageReceivedStop(message mqtt.Message) {
	speaker.Clear()

	a.lastPlayedUid = 0
	a.sendStatusMessage(mqtt.Info, "stopped")
}

func (a *Audio) OnMessageReceivedVolume(message mqtt.Message) {
	if a.volume != nil {
		volumeMessage := &mqtt.VolumeChangeSubscriptionMessage{}

		message.ToStruct(volumeMessage)

		speaker.Lock()
		a.volume.Volume += volumeMessage.VolumeOffset
		speaker.Unlock()

		a.sendStatusMessage(mqtt.Info, "Volume changed by ", volumeMessage.VolumeOffset)
	} else {
		a.sendStatusMessage(mqtt.Warn, "Volume not changed, because nothing is played...")
	}
}

func NewAudio(statusMessage func(messageType mqtt.StatusType, message ...any), playDoneMessage func(id uint)) *Audio {
	return &Audio{
		sendStatusMessage:   statusMessage,
		sendPlayDoneMessage: playDoneMessage,
	}
}
