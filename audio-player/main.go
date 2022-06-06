package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
	"os"
	"pmoscode.diy-media-player-box-backend/audio"
	"time"
)

func main() {
	port := os.Getenv("AUDIO_PLAYER_PORT")

	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	const sampleRate = beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	router.POST("/audio", audio.Play)
	router.PATCH("/audio", audio.SwitchPlayState)
	router.DELETE("/audio", audio.Stop)

	router.Run(port)
}
