package main

import (
	"flag"
	"log"
	"os/exec"
	"rfid-reader/mqtt"
	"rfid-reader/rfid"
)

var mqttClient *mqtt.Client

type CliOptions struct {
	mqttBrokerIp *string
	mqttClientId *string
	mockCardId   *string
}

type Module interface {
	Run()
}

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "rfid-reader", "Client id for Mqtt connection")
	mockCardId := flag.String("mock-card-id", "123456", "Only used when in mock mode")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return CliOptions{
		mqttBrokerIp: mqttBrokerIp,
		mqttClientId: mqttClientId,
		mockCardId:   mockCardId,
	}
}

func main() {
	cliOptions := getCliOptions()

	mqttClient = mqtt.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	mqttClient.Connect()

	var rfidClient Module

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Not on Raspi... Switching to Mock mode...")
		rfidClient = rfid.NewMock(cliOptions.mockCardId, mqttClient)
	} else {
		log.Println("On Raspi... Switching to Rfid mode...")
		rfidClient = rfid.NewRfid(sendCardIdMessage, sendStatusMessage)
	}
	rfidClient.Run()
}

func sendStatusMessage(statusMmessage string) {
	message := &rfid.StatusPublishMessage{
		Status: statusMmessage,
	}

	mqttClient.SendMessage(&mqtt.Message{
		Topic: "/status/RfidReader",
		Value: message,
	})
}

func sendCardIdMessage(cardId string) {
	message := &rfid.CardIdPublishMessage{
		CardId: cardId,
	}

	mqttClient.SendMessage(&mqtt.Message{
		Topic: "/rfidReader/cardId",
		Value: message,
	})
}
