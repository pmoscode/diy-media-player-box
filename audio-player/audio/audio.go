package audio

import (
	"audio-player/mqtt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
)

const DefaultSampleRate = 44100

type Audio struct {
	control           *beep.Ctrl
	lastPlayedUid     string
	sendStatusMessage func(message string)
}

func (a *Audio) checkLastPlayedUidChanged(body *TracksSubscriptionMessage) bool {
	if a.lastPlayedUid != body.Id {
		a.lastPlayedUid = body.Id

		return true
	}

	return false
}

func (a *Audio) OnMessageReceivedPlay(message mqtt.Message) {
	body := TracksSubscriptionMessage{}

	message.ToStruct(&body)

	uidChanged := a.checkLastPlayedUidChanged(&body)
	log.Println("uid changed: ", uidChanged)

	if uidChanged {
		speaker.Clear()

		var samples []beep.Streamer

		for _, trackPath := range body.TrackList {
			f, err := os.Open(trackPath)
			if err != nil {
				log.Fatal(err)
			}

			streamer, format, err := mp3.Decode(f)
			if err != nil {
				log.Fatal(err)
			}

			var stream beep.Streamer

			if DefaultSampleRate != format.SampleRate {
				const sampleRate = beep.SampleRate(DefaultSampleRate)
				stream = beep.Resample(1, format.SampleRate, sampleRate, streamer)
				log.Println("Need to resample: ", trackPath)
			} else {
				stream = streamer
				log.Println("No need to resample: ", trackPath)
			}

			samples = append(samples, stream)
		}

		samples = append(samples, beep.Callback(func() {
			a.lastPlayedUid = "-1"
			log.Println("status: ", "stopped")
			a.sendStatusMessage("stopped")
		}))

		if len(samples) > 0 {
			sequence := beep.Seq(samples...)
			a.control = &beep.Ctrl{Streamer: sequence, Paused: false}

			speaker.Play(a.control)

			log.Println("status: ", "playing")
			a.sendStatusMessage("playing")
		} else {
			log.Println("status: ", "no tracks")
			a.sendStatusMessage("no tracks")
		}
	} else {
		if a.control.Paused {
			a.OnMessageReceivedSwitch(message)
		} else {
			log.Println("status: ", "untouched")
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

		log.Println("status: ", status)
		a.sendStatusMessage(status)
	} else {
		log.Println("status: ", "no audio stream")
		a.sendStatusMessage("no audio stream")
	}
}

func (a *Audio) OnMessageReceivedStop(message mqtt.Message) {
	speaker.Clear()

	a.lastPlayedUid = "-1"
	log.Println("status: ", "stopped")
	a.sendStatusMessage("stopped")
}

func NewAudio(statusMessage func(statusMessage string)) *Audio {
	return &Audio{sendStatusMessage: statusMessage}
}
