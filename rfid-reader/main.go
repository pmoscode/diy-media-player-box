package main

import (
	"fmt"
	"gitlab.com/pmoscodegrp/common/heartbeat"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"gitlab.com/pmoscodegrp/common/yamlconfig"
	"log"
	"os/exec"
	"rfid-reader/mqtt"
	"rfid-reader/rfid"
	"time"
)

var mqttClient *mqtt2.Client

type Module interface {
	Run()
}

var config Config

func main() {
	err := yamlconfig.LoadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("Could not load config file")
	}

	mqttClient = mqtt2.CreateClient(config.MqttBroker.Host, 1883, config.RfidReader.MqttClientId)
	err = mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	heartBeat := heartbeat.New(10*time.Second, sendHeartbeat)
	heartBeat.Run()

	var rfidClient Module

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	_, err = cmd.CombinedOutput()
	if err != nil {
		sendStatusMessage(mqtt2.Info, "Not on Raspi... Switching to Mock mode...")
		rfidClient = rfid.NewMock(config.RfidReader.MockCardId, sendCardIdMessage, sendStatusMessage)
	} else {
		sendStatusMessage(mqtt2.Info, "On Raspi... Switching to Rfid mode...")
		rfidClient = rfid.NewRfid(config.RfidReader.RemoveThreshold, sendCardIdMessage, sendStatusMessage)
	}
	rfidClient.Run()
}

func sendStatusMessage(messageType mqtt2.StatusType, message ...any) {
	messageTxt := fmt.Sprint(message...)

	mqttMessage := &mqtt2.StatusPublishMessage{
		Type:      messageType,
		Status:    messageTxt,
		Timestamp: time.Now(),
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/status/rfid-reader",
		Value: mqttMessage,
	})

	if config.RfidReader.LogStatusToConsole {
		log.Println(messageType, ": ", messageTxt)
	}
}

func sendCardIdMessage(cardId string) {
	message := &mqtt.CardIdPublishMessage{
		CardId: cardId,
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/rfid-reader/cardId",
		Value: message,
	})
}

func sendHeartbeat() {
	mqttClient.Publish(&mqtt2.Message{
		Topic: "/heartbeat/rfid-reader",
		Value: &heartbeat.PublishMessage{
			Alive: true,
		},
	})
}
