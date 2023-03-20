package main

import (
	"fmt"
	"gitlab.com/pmoscodegrp/common/heartbeat"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"gitlab.com/pmoscodegrp/common/yamlconfig"
	config2 "io-controller/config"
	"io-controller/io"
	"io-controller/mqtt"
	"log"
	"os/exec"
	"time"
)

var mqttClient *mqtt2.Client
var config config2.Config

type Module interface {
	Run()
}

func main() {
	err := yamlconfig.LoadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("Could not load config file")
	}

	mqttClient = mqtt2.CreateClient(config.MqttBroker.Host, config.MqttBroker.Port, config.IoController.MqttClientId)
	err = mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	heartBeat := heartbeat.New(10*time.Second, sendHeartbeat)
	heartBeat.Run()

	var ioClient Module

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	_, err = cmd.CombinedOutput()
	if err != nil {
		sendStatusMessage(mqtt2.Info, "Not on Raspi... Switching to Mock mode...")
		ioClient = io.NewMockOI(config.IoController.MockVolumeOffset, sendVolumeChangeMessage, sendStatusMessage)
	} else {
		sendStatusMessage(mqtt2.Info, "On Raspi... Switching to IO mode...")
		ioClient = io.NewOI(config.IoController, sendVolumeChangeMessage, sendTrackChangeMessage, sendStatusMessage)
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

	if config.IoController.LogStatusToConsole {
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

func sendHeartbeat() {
	mqttClient.Publish(&mqtt2.Message{
		Topic: "/heartbeat/io-controller",
		Value: &heartbeat.PublishMessage{
			Alive: true,
		},
	})
}
