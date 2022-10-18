package main

import (
	"audio-player/audio"
	"audio-player/mqtt"
	"flag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"log"
	"time"
)

var mqttClient *mqtt.Client

type CliOptions struct {
	mqttBrokerIp     *string
	mqttClientId     *string
	sampleRateFactor *int
}

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "audio-player", "Client id for Mqtt connection")
	sampleRateFactor := flag.Int("sample-rate-factor", 10, "Buffer size of audio player")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return CliOptions{
		mqttBrokerIp:     mqttBrokerIp,
		mqttClientId:     mqttClientId,
		sampleRateFactor: sampleRateFactor,
	}
}

func main() {
	cliOptions := getCliOptions()

	mqttClient = mqtt.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	mqttClient.Connect()

	const sampleRate = beep.SampleRate(audio.DefaultSampleRate)
	speaker.Init(sampleRate, sampleRate.N(time.Second/time.Duration(*cliOptions.sampleRateFactor)))

	audioClient := audio.NewAudio(sendStatusMessage, sendPlayDoneMessage)

	mqttClient.Subscribe("/audioPlayer/play", audioClient.OnMessageReceivedPlay)
	mqttClient.Subscribe("/audioPlayer/switch", audioClient.OnMessageReceivedSwitch)
	mqttClient.Subscribe("/audioPlayer/stop", audioClient.OnMessageReceivedStop)
	mqttClient.LoopForever()
}

func sendStatusMessage(message string) {
	publishMessage := &audio.StatusPublishMessage{
		Status: message,
	}

	mqttClient.SendMessage(&mqtt.Message{
		Topic: "/status/audioPlayer",
		Value: publishMessage,
	})
}

func sendPlayDoneMessage(id uint) {
	publishMessage := &audio.PlayDonePublishMessage{
		Uid: id,
	}

	mqttClient.SendMessage(&mqtt.Message{
		Topic: "/audioPlayer/done",
		Value: publishMessage,
	})
}
