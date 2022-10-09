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
	mqttBrokerIp     *string
	mqttClientId     *string
	sampleRateFactor *int
}

type Module interface {
	Run()
}

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "rfid-reader", "Client id for Mqtt connection")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return CliOptions{
		mqttBrokerIp: mqttBrokerIp,
		mqttClientId: mqttClientId,
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
		rfidClient = rfid.NewMock(mqttClient)
	} else {
		log.Println("On Raspi... Switching to Rfid mode...")
		rfidClient = rfid.NewRfid(sendCardIdMessage, sendStatusMessage)
	}
	rfidClient.Run()
}

func sendStatusMessage(message string) {
	mqttClient.SendMessage(&mqtt.Message{
		Topic: "/status/RfidReader",
		Value: "{\"status\": \"" + message + "\"}",
	})
}

func sendCardIdMessage(cardId string) {
	mqttClient.SendMessage(&mqtt.Message{
		Topic: "/rfidReader/cardId",
		Value: "{\"cardId\": \"" + cardId + "\"}",
	})
}
