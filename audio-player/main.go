package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
	"pmoscode.diy-media-player-box-backend/audio"
	"pmoscode.diy-media-player-box-backend/health"
	"time"
)

func main() {
	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	const sampleRate = beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	r.POST("/audio", audio.Play)
	r.PATCH("/audio", audio.Play)
	r.DELETE("/audio", audio.Stop)
	r.GET("/health", health.Health)

	r.Run()
}
