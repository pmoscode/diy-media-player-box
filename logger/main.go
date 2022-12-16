package main

import (
	"flag"
	"gitlab.com/pmoscodegrp/common/mqtt"
	"log"
	"logger/logs"
)

var mqttClient *mqtt.Client

type CliOptions struct {
	mqttBrokerIp                *string
	mqttClientId                *string
	mqttSubscriptionTopic       *string
	fileName                    *string
	logRotationPeriodAfterBytes *int
}

func getCliOptions() CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "logger", "Client id for Mqtt connection")
	mqttSubscriptionTopic := flag.String("mqtt-sub-topic", "/status/#", "Topic to subscribe to")
	fileName := flag.String("filename", "logs/music.log", "Defines the filename of the log file")
	logRotationPeriodAfterBytes := flag.Int("file-size", 10000000, "Maximum file size of the log in bytes (max == 2GB). Then log-file will be rotated.")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return CliOptions{
		mqttBrokerIp:                mqttBrokerIp,
		mqttClientId:                mqttClientId,
		fileName:                    fileName,
		logRotationPeriodAfterBytes: logRotationPeriodAfterBytes,
		mqttSubscriptionTopic:       mqttSubscriptionTopic,
	}
}
func main() {
	cliOptions := getCliOptions()

	mqttClient = mqtt.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	err := mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	logger := logs.New(*cliOptions.fileName, *cliOptions.logRotationPeriodAfterBytes)

	log.Printf("Subscribing on '%s'\n", *cliOptions.mqttSubscriptionTopic)

	mqttClient.Subscribe(*cliOptions.mqttSubscriptionTopic, logger.Log)
	mqttClient.LoopForever()
}
