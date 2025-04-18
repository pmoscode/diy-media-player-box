package main

import (
	"audio-player/audio"
	"audio-player/mqtt"
	"fmt"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/pmoscode/go-common/heartbeat"
	mqtt2 "github.com/pmoscode/go-common/mqtt"
	"github.com/pmoscode/go-common/shutdown"
	"github.com/pmoscode/go-common/yamlconfig"
	"log"
	"time"
)

var mqttClient *mqtt2.Client

var config Config

func main() {
	defer shutdown.ExitOnPanic()

	err := yamlconfig.LoadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("Could not load config file")
	}

	mqttClient = mqtt2.NewClient(mqtt2.WithBroker(config.MqttBroker.Host, 1883),
		mqtt2.WithClientId(config.AudioPlayer.MqttClientId),
		mqtt2.WithOrderMatters(false))
	err = mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}
	defer mqttClient.Disconnect()

	const sampleRate = beep.SampleRate(audio.DefaultSampleRate)
	err = speaker.Init(sampleRate, sampleRate.N(time.Duration(config.AudioPlayer.SampleRateFactor)*time.Millisecond))
	if err != nil {
		log.Fatal(err)
	}

	heartBeat := heartbeat.New(10*time.Second, sendHeartbeat)
	heartBeat.Run()

	audioClient := audio.NewAudio(sendStatusMessage, sendPlayDoneMessage)

	mqttClient.Subscribe("/controller/play", audioClient.OnMessageReceivedPlay)
	mqttClient.Subscribe("/controller/pause", audioClient.OnMessageReceivedPause)
	mqttClient.Subscribe("/controller/resume", audioClient.OnMessageReceivedResume)
	mqttClient.Subscribe("/controller/stop", audioClient.OnMessageReceivedStop)
	mqttClient.Subscribe("/io-controller/volume", audioClient.OnMessageReceivedVolume)
	mqttClient.LoopForever()
}

func sendStatusMessage(messageType mqtt2.StatusType, message ...any) {
	messageTxt := fmt.Sprint(message...)

	mqttMessage := &mqtt2.StatusPublishMessage{
		Type:      messageType,
		Status:    messageTxt,
		Timestamp: time.Now(),
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/status/audio-player",
		Value: mqttMessage,
	})

	if config.AudioPlayer.LogStatusToConsole {
		log.Println(messageType, ": ", messageTxt)
	}
}

func sendPlayDoneMessage(id uint) {
	publishMessage := &mqtt.PlayDonePublishMessage{
		Uid: id,
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/audio-player/done",
		Value: publishMessage,
	})
}

func sendHeartbeat() {
	mqttClient.Publish(&mqtt2.Message{
		Topic: "/heartbeat/audio-player",
		Value: &mqtt2.StatusPublishMessage{
			Status: "online",
			Type:   mqtt2.Info,
		},
	})
}
