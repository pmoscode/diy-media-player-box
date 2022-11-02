package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"rfid-reader/mqtt"
	"rfid-reader/rfid"
)

var mqttClient *mqtt.Client

type CliOptions struct {
	mqttBrokerIp       *string
	mqttClientId       *string
	mockCardId         *string
	logStatusToConsole *bool
}

type Module interface {
	Run()
}

var cliOptions CliOptions

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "rfid-reader", "Client id for Mqtt connection")
	mockCardId := flag.String("mock-card-id", "123456", "Only used when in mock mode")
	logStatusToConsole := flag.Bool("log-console", false, "Log messages also to current std console")
	flag.Parse()

	return CliOptions{
		mqttBrokerIp:       mqttBrokerIp,
		mqttClientId:       mqttClientId,
		mockCardId:         mockCardId,
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

	var rfidClient Module

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	_, err = cmd.CombinedOutput()
	if err != nil {
		sendStatusMessage(mqtt.Info, "Not on Raspi... Switching to Mock mode...")
		rfidClient = rfid.NewMock(cliOptions.mockCardId, sendCardIdMessage, sendStatusMessage)
	} else {
		sendStatusMessage(mqtt.Info, "On Raspi... Switching to Rfid mode...")
		rfidClient = rfid.NewRfid(sendCardIdMessage, sendStatusMessage)
	}
	rfidClient.Run()
}

func sendStatusMessage(messageType mqtt.StatusType, message ...any) {
	messageTxt := fmt.Sprint(message...)

	mqttMessage := &mqtt.StatusPublishMessage{
		Type:   messageType,
		Status: messageTxt,
	}

	mqttClient.Publish(&mqtt.Message{
		Topic: "/status/rfid-reader",
		Value: mqttMessage,
	})

	if *cliOptions.logStatusToConsole {
		log.Println(messageType, ": ", messageTxt)
	}
}

func sendCardIdMessage(cardId string) {
	message := &mqtt.CardIdPublishMessage{
		CardId: cardId,
	}

	mqttClient.Publish(&mqtt.Message{
		Topic: "/rfid-reader/cardId",
		Value: message,
	})
}