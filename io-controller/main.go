package main

import (
	"fmt"
	"io-controller/cli"
	"io-controller/io"
	"io-controller/mqtt"
	"log"
	"os/exec"
)

var mqttClient *mqtt.Client
var cliOptions *cli.Options

type Module interface {
	Run()
}

func main() {
	// Pins 20 (38), 16 (36), 13 (33), 12 (32)

	cliOptions = cli.GetCliOptions()

	mqttClient = mqtt.CreateClient(*cliOptions.MqttBrokerIp, 1883, *cliOptions.MqttClientId)
	err := mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	var ioClient Module

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	_, err = cmd.CombinedOutput()
	if err != nil {
		sendStatusMessage(mqtt.Info, "Not on Raspi... Switching to Mock mode...")
		ioClient = io.NewMockOI(*cliOptions.MockVolumeOffset, sendVolumeChangeMessage, sendStatusMessage)
	} else {
		sendStatusMessage(mqtt.Info, "On Raspi... Switching to IO mode...")
		ioClient = io.NewOI(cliOptions, sendVolumeChangeMessage, sendTrackChangeMessage, sendStatusMessage)
	}

	ioClient.Run()
	mqttClient.Disconnect()
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

	if *cliOptions.LogStatusToConsole {
		log.Println(messageType, ": ", messageTxt)
	}
}

func sendVolumeChangeMessage(volumeOffset float64) {
	publishMessage := &mqtt.VolumeChangePublishMessage{
		VolumeOffset: volumeOffset,
	}

	mqttClient.Publish(&mqtt.Message{
		Topic: "/audio-player/volume",
		Value: publishMessage,
	})
}

func sendTrackChangeMessage(direction int) {
	publishMessage := &mqtt.TrackChangePublishMessage{
		Direction: direction,
	}

	mqttClient.Publish(&mqtt.Message{
		Topic: "/audio-player/track",
		Value: publishMessage,
	})
}
