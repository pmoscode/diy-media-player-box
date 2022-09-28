package main

import (
	"audio-player/audio"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

func getParameterFromEnvironment() (string, int) {
	address := os.Getenv("AUDIO_PLAYER_BIND_ADDRESS")
	sampleRateFactorEnv := os.Getenv("AUDIO_PLAYER_SAMPLE_RATE_FACTOR")

	if address == "" {
		log.Println("Env variable 'AUDIO_PLAYER_BIND_ADDRESS' is not set properly (ex.: export AUDIO_PLAYER_BIND_ADDRESS=localhost:8080)")
		os.Exit(1)
	}

	if sampleRateFactorEnv == "" {
		log.Println("Env variable 'AUDIO_PLAYER_SAMPLE_RATE_FACTOR' is not set (ex.: export AUDIO_PLAYER_SAMPLE_RATE_FACTOR=1) - default: 10")
		sampleRateFactorEnv = "10"
	}

	sampleRateFactor, _ := strconv.Atoi(sampleRateFactorEnv)
	log.Println("Configured bind address: ", address)
	log.Println("Configured SampleFactor: ", sampleRateFactor)

	return address, sampleRateFactor
}

func main() {
	address, sampleRateFactor := getParameterFromEnvironment()

	const sampleRate = beep.SampleRate(audio.SampleRate)
	speaker.Init(sampleRate, sampleRate.N(time.Second/time.Duration(sampleRateFactor)))

	router := gin.Default()
	router.POST("/audio", audio.Play)
	router.PATCH("/audio", audio.SwitchPlayState)
	router.DELETE("/audio", audio.Stop)
	router.Run(address)
}
