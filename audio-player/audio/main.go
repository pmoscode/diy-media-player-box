package audio

import (
	"audio-player/mqtt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"log"
	"os"
)

const DefaultSampleRate = 44100

type Audio struct {
	control             *beep.Ctrl
	volume              *effects.Volume
	lastPlayedUid       uint
	sendStatusMessage   func(messageType mqtt2.StatusType, message ...any)
	sendPlayDoneMessage func(id uint)
	currentVolume       float64
}

func (a *Audio) OnMessageReceivedPlay(message mqtt2.Message) {
	body := mqtt.TracksSubscriptionMessage{}

	message.ToStruct(&body)
	a.sendStatusMessage(mqtt2.Info, "Got audio book: '", body.Id, "' and track list of size: ", len(body.TrackList))

	speaker.Clear()

	if len(body.TrackList) == 0 {
		a.sendStatusMessage(mqtt2.Warn, "no tracks")

		return
	}

	var samples []beep.Streamer

	for _, trackPath := range body.TrackList {
		f, err := os.Open(trackPath)
		if err != nil {
			a.sendStatusMessage(mqtt2.Error, "Could not open '"+trackPath+"'... DYING!!!")
			log.Fatal(err)
		}

		streamer, format, err := mp3.Decode(f)
		if err != nil {
			a.sendStatusMessage(mqtt2.Error, "Could not decode '"+trackPath+"'... DYING!!!")
			log.Fatal(err)
		}

		var stream beep.Streamer

		if DefaultSampleRate != format.SampleRate {
			const sampleRate = beep.SampleRate(DefaultSampleRate)
			stream = beep.Resample(1, format.SampleRate, sampleRate, streamer)
			a.sendStatusMessage(mqtt2.Info, "Need to resample: '"+trackPath+"'...")
		} else {
			stream = streamer
			a.sendStatusMessage(mqtt2.Info, "No need to resample: '"+trackPath+"'...")
		}

		samples = append(samples, stream)
	}

	samples = append(samples, beep.Callback(func() {
		a.lastPlayedUid = 0
		a.sendStatusMessage(mqtt2.Info, "stopped")
		a.sendPlayDoneMessage(body.Id)
		a.control = nil
		a.volume = nil
		// speaker.Clear()
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
			Volume:   a.currentVolume,
			Silent:   false,
		}

		speaker.Play(a.volume)

		a.sendStatusMessage(mqtt2.Info, "playing")
	}
}

func (a *Audio) OnMessageReceivedPause(message mqtt2.Message) {
	if a.control != nil {
		speaker.Lock()
		a.control.Paused = true
		speaker.Unlock()

		a.sendStatusMessage(mqtt2.Info, "paused")
	} else {
		a.sendStatusMessage(mqtt2.Info, "no audio stream to pause...")
	}
}

func (a *Audio) OnMessageReceivedResume(message mqtt2.Message) {
	if a.control != nil {
		speaker.Lock()
		a.control.Paused = false
		speaker.Unlock()

		a.sendStatusMessage(mqtt2.Info, "continuing")
	} else {
		a.sendStatusMessage(mqtt2.Info, "no audio stream to continue...")
	}
}

func (a *Audio) OnMessageReceivedStop(message mqtt2.Message) {
	//speaker.Clear()

	a.lastPlayedUid = 0
	a.sendStatusMessage(mqtt2.Info, "stopped")
}

func (a *Audio) OnMessageReceivedVolume(message mqtt2.Message) {
	if a.volume != nil {
		volumeMessage := &mqtt.VolumeChangeSubscriptionMessage{}

		message.ToStruct(volumeMessage)

		speaker.Lock()
		a.volume.Volume += volumeMessage.VolumeOffset
		a.currentVolume = a.volume.Volume
		speaker.Unlock()

		a.sendStatusMessage(mqtt2.Info, "Volume changed by ", volumeMessage.VolumeOffset)
	} else {
		a.sendStatusMessage(mqtt2.Warn, "Volume not changed, because nothing is played...")
	}
}

func NewAudio(statusMessage func(messageType mqtt2.StatusType, message ...any), playDoneMessage func(id uint)) *Audio {
	return &Audio{
		sendStatusMessage:   statusMessage,
		sendPlayDoneMessage: playDoneMessage,
		currentVolume:       0.0,
	}
}
