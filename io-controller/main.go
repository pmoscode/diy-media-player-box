package main

import (
	"fmt"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"io-controller/cli"
	"io-controller/io"
	"io-controller/mqtt"
	"log"
	"os/exec"
	"time"
)

var mqttClient *mqtt2.Client
var cliOptions *cli.Options

type Module interface {
	Run()
}

func main() {
	cliOptions = cli.GetCliOptions()

	mqttClient = mqtt2.CreateClient(*cliOptions.MqttBrokerIp, 1883, *cliOptions.MqttClientId)
	err := mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	var ioClient Module

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	_, err = cmd.CombinedOutput()
	if err != nil {
		sendStatusMessage(mqtt2.Info, "Not on Raspi... Switching to Mock mode...")
		ioClient = io.NewMockOI(*cliOptions.MockVolumeOffset, sendVolumeChangeMessage, sendStatusMessage)
	} else {
		sendStatusMessage(mqtt2.Info, "On Raspi... Switching to IO mode...")
		ioClient = io.NewOI(cliOptions, sendVolumeChangeMessage, sendTrackChangeMessage, sendStatusMessage)
	}

	ioClient.Run()
	mqttClient.Disconnect()
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

	if *cliOptions.LogStatusToConsole {
		log.Println(messageType, ": ", messageTxt)
	}
}

func sendVolumeChangeMessage(volumeOffset float64) {
	publishMessage := &mqtt.VolumeChangePublishMessage{
		VolumeOffset: volumeOffset,
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/io-controller/volume",
		Value: publishMessage,
	})
}

func sendTrackChangeMessage(direction int) {
	publishMessage := &mqtt.TrackChangePublishMessage{
		Direction: direction,
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/io-controller/track",
		Value: publishMessage,
	})
}
