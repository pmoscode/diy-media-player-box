package main

import (
	mqtt2 "github.com/pmoscode/go-common/mqtt"
	"github.com/pmoscode/go-common/yamlconfig"
	"log"
	"monitor/monit"
)

var mqttClient *mqtt2.Client

var config Config

func main() {
	err := yamlconfig.LoadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("Could not load config file")
	}

	mqttClient = mqtt2.NewClient(mqtt2.WithBroker(config.MqttBroker.Host, 1883),
		mqtt2.WithClientId(config.Monitor.MqttClientId),
		mqtt2.WithOrderMatters(false))
	err = mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}
	defer mqttClient.Disconnect()

	mqttClient.Subscribe("/heartbeat/#", monit.OnMessageReceivedHeartbeat)

	monit.RunMonitor(config.Monitor.ProcessNames)

	mqttClient.LoopForever()
}
