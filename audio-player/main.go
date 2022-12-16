package main

import (
	"audio-player/audio"
	"audio-player/mqtt"
	"flag"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"gitlab.com/pmoscodegrp/common/heartbeat"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"log"
	"time"
)

var mqttClient *mqtt2.Client

type CliOptions struct {
	mqttBrokerIp       *string
	mqttClientId       *string
	sampleRateFactor   *int
	logStatusToConsole *bool
}

var cliOptions CliOptions

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "audio-player", "Client id for Mqtt connection")
	sampleRateFactor := flag.Int("buffer-sample-rate", 400, "Defines buffer size of audio player # in milliseconds")
	logStatusToConsole := flag.Bool("log-console", false, "Log messages also to current std console")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return CliOptions{
		mqttBrokerIp:       mqttBrokerIp,
		mqttClientId:       mqttClientId,
		sampleRateFactor:   sampleRateFactor,
		logStatusToConsole: logStatusToConsole,
	}
}

func main() {
	cliOptions = getCliOptions()

	mqttClient = mqtt2.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	err := mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	const sampleRate = beep.SampleRate(audio.DefaultSampleRate)
	err = speaker.Init(sampleRate, sampleRate.N(time.Duration(*cliOptions.sampleRateFactor)*time.Millisecond))
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

	if *cliOptions.logStatusToConsole {
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
		Value: &heartbeat.PublishMessage{
			Alive: true,
		},
	})
}
