package main

import (
	"flag"
	"fmt"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"log"
	"os/exec"
	"rfid-reader/mqtt"
	"rfid-reader/rfid"
	"time"
)

var mqttClient *mqtt2.Client

type CliOptions struct {
	mqttBrokerIp       *string
	mqttClientId       *string
	mockCardId         *string
	logStatusToConsole *bool
	removeThreshold    *int
}

type Module interface {
	Run()
}

var cliOptions CliOptions

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "rfid-reader", "Client id for Mqtt connection")
	mockCardId := flag.String("mock-card-id", "123456", "Only used when in mock mode")
	removeThreshold := flag.Int("remove-threshold", 2, "How many checks for removed card until it will be noticed as 'card removed'")
	logStatusToConsole := flag.Bool("log-console", false, "Log messages also to current std console")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return CliOptions{
		mqttBrokerIp:       mqttBrokerIp,
		mqttClientId:       mqttClientId,
		mockCardId:         mockCardId,
		logStatusToConsole: logStatusToConsole,
		removeThreshold:    removeThreshold,
	}
}

func main() {
	cliOptions = getCliOptions()

	mqttClient = mqtt2.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	err := mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	var rfidClient Module

	cmd := exec.Command("cat", "/sys/firmware/devicetree/base/serial-number")
	_, err = cmd.CombinedOutput()
	if err != nil {
		sendStatusMessage(mqtt2.Info, "Not on Raspi... Switching to Mock mode...")
		rfidClient = rfid.NewMock(cliOptions.mockCardId, sendCardIdMessage, sendStatusMessage)
	} else {
		sendStatusMessage(mqtt2.Info, "On Raspi... Switching to Rfid mode...")
		rfidClient = rfid.NewRfid(cliOptions.removeThreshold, sendCardIdMessage, sendStatusMessage)
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

	if *cliOptions.logStatusToConsole {
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
