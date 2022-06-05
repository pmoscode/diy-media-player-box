package audio

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var control *beep.Ctrl

func Play(c *gin.Context) {
	speaker.Clear()

	f, err := os.Open("/home/peter/Arbeit/GIT/GitLab/GoLang/diy-media-player-box-backend/test/test.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	const sampleRate = beep.SampleRate(44100)
	resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)

	fmt.Println("Before control")
	control = &beep.Ctrl{Streamer: resampled, Paused: false}
	fmt.Println("Before play")
	speaker.Play(control)

	c.JSON(http.StatusOK, gin.H{"status": "playing..."})
}

func SwitchPlayState(c *gin.Context) {
	speaker.Lock()
	control.Paused = !control.Paused
	speaker.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "switched playing state..."})
}

func Stop(c *gin.Context) {
	speaker.Lock()
	control.Paused = true
	speaker.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "stopped playing..."})
}
