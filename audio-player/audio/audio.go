package audio

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

const SampleRate = 44100

var control *beep.Ctrl

var lastPlayedUid = "-1"

func Play(context *gin.Context) {

	var body audioRequestInput

	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("body: ", body)

	uidChanged := checkLastPlayedUidChanged(&body)
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

			if SampleRate != format.SampleRate {
				const sampleRate = beep.SampleRate(SampleRate)
				stream = beep.Resample(1, format.SampleRate, sampleRate, streamer)
				log.Println("Need to resample: ", trackPath)
			} else {
				stream = streamer
				log.Println("No need to resample: ", trackPath)
			}

			samples = append(samples, stream)
		}

		samples = append(samples, beep.Callback(func() {
			lastPlayedUid = "-1"
		}))

		if len(samples) > 0 {
			sequence := beep.Seq(samples...)
			control = &beep.Ctrl{Streamer: sequence, Paused: false}

			speaker.Play(control)

			context.JSON(http.StatusOK, gin.H{"status": "playing"})
		} else {
			context.JSON(http.StatusOK, gin.H{"status": "no tracks"})
		}
	} else {
		if control.Paused {
			SwitchPlayState(context)
		} else {
			context.JSON(http.StatusOK, gin.H{"status": "untouched"})
		}
	}
}

func checkLastPlayedUidChanged(body *audioRequestInput) bool {
	if lastPlayedUid != body.Uid {
		lastPlayedUid = body.Uid

		return true
	}

	return false
}

func SwitchPlayState(context *gin.Context) {
	speaker.Lock()
	control.Paused = !control.Paused
	speaker.Unlock()

	status := "paused"
	if !control.Paused {
		status = "continuing"
	}

	context.JSON(http.StatusOK, gin.H{"status": status})
}

func Stop(context *gin.Context) {
	speaker.Clear()

	lastPlayedUid = "-1"

	context.JSON(http.StatusOK, gin.H{"status": "stopped"})
}
