package main

import (
	"flag"
	"fmt"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"log"
	"monitor/monit"
	"time"
)

var mqttClient *mqtt2.Client

type CliOptions struct {
	mqttBrokerIp       *string
	mqttClientId       *string
	logStatusToConsole *bool
	processNames       *string
}

var cliOptions CliOptions

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "monitor", "Client id for Mqtt connection")
	processNames := flag.String("process-names", "audio-player,controller,io-controller,logger,rfid-reader", "define the processes to watch for")
	logStatusToConsole := flag.Bool("log-console", false, "Log messages also to current std console")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return CliOptions{
		mqttBrokerIp:       mqttBrokerIp,
		mqttClientId:       mqttClientId,
		logStatusToConsole: logStatusToConsole,
		processNames:       processNames,
	}
}

func main() {
	cliOptions = getCliOptions()

	mqttClient = mqtt2.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	err := mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}
	defer mqttClient.Disconnect()

	mqttClient.Subscribe("/heartbeat/#", monit.OnMessageReceivedHeartbeat)

	monit.RunMonitor(cliOptions.processNames)

	mqttClient.LoopForever()
}

func sendStatusMessage(messageType mqtt2.StatusType, message ...any) {
	messageTxt := fmt.Sprint(message...)

	mqttMessage := &mqtt2.StatusPublishMessage{
		Type:      messageType,
		Status:    messageTxt,
		Timestamp: time.Now(),
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/status/monitor",
		Value: mqttMessage,
	})

	if *cliOptions.logStatusToConsole {
		log.Println(messageType, ": ", messageTxt)
	}
}
