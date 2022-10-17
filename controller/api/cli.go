package api

import (
	"flag"
	"log"
)

type CliOptions struct {
	mqttBrokerIp *string
	mqttClientId *string
}

func getCliOptions() *CliOptions {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "controller", "Client id for Mqtt connection")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return &CliOptions{
		mqttBrokerIp: mqttBrokerIp,
		mqttClientId: mqttClientId,
	}
}

var cliOptions = getCliOptions()
