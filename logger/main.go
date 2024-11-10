package main

import (
	"github.com/pmoscode/go-common/heartbeat"
	"github.com/pmoscode/go-common/mqtt"
	"github.com/pmoscode/go-common/shutdown"
	"github.com/pmoscode/go-common/yamlconfig"
	"log"
	"logger/logs"
	"time"
)

var mqttClient *mqtt.Client

var config Config

func main() {
	defer shutdown.ExitOnPanic()

	err := yamlconfig.LoadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("Could not load config file")
	}

	mqttClient = mqtt.NewClient(mqtt.WithBroker(config.MqttBroker.Host, 1883),
		mqtt.WithClientId(config.Logger.MqttClientId),
		mqtt.WithOrderMatters(false))
	err = mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}

	heartBeat := heartbeat.New(10*time.Second, sendHeartbeat)
	heartBeat.Run()

	logger := logs.New(config.Logger.FileName, config.Logger.LogRotationPeriodAfterBytes)

	log.Printf("Subscribing on '%s'\n", config.Logger.MqttSubscriptionTopic)

	mqttClient.Subscribe(config.Logger.MqttSubscriptionTopic, logger.Log)
	mqttClient.LoopForever()
}

func sendHeartbeat() {
	mqttClient.Publish(&mqtt.Message{
		Topic: "/heartbeat/logger",
		Value: &mqtt.StatusPublishMessage{
			Status: "online",
			Type:   mqtt.Info,
		},
	})
}
