package main

import (
	"audio-player/audio"
	"audio-player/mqtt"
	"flag"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"log"
	"time"
)

var mqttClient *mqtt.Client

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
	sampleRateFactor := flag.Int("sample-rate-factor", 10, "Buffer size of audio player")
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

	mqttClient = mqtt.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	err := mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	const sampleRate = beep.SampleRate(audio.DefaultSampleRate)
	speaker.Init(sampleRate, sampleRate.N(time.Second/time.Duration(*cliOptions.sampleRateFactor)))

	audioClient := audio.NewAudio(sendStatusMessage, sendPlayDoneMessage)

	mqttClient.Subscribe("/audioPlayer/play", audioClient.OnMessageReceivedPlay)
	mqttClient.Subscribe("/audioPlayer/switch", audioClient.OnMessageReceivedSwitch)
	mqttClient.Subscribe("/audioPlayer/stop", audioClient.OnMessageReceivedStop)
	mqttClient.Subscribe("/audioPlayer/volume", audioClient.OnMessageReceivedVolume)
	mqttClient.LoopForever()
}

func sendStatusMessage(messageType mqtt.StatusType, message ...any) {
	messageTxt := fmt.Sprint(message...)

	mqttMessage := &mqtt.StatusPublishMessage{
		Type:   messageType,
		Status: messageTxt,
	}

	mqttClient.Publish(&mqtt.Message{
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

	mqttClient.Publish(&mqtt.Message{
		Topic: "/audioPlayer/done",
		Value: publishMessage,
	})
}
