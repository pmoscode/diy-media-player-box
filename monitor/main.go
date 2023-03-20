package main

import (
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"gitlab.com/pmoscodegrp/common/yamlconfig"
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

	mqttClient = mqtt2.CreateClient(config.MqttBroker.Host, config.MqttBroker.Port, config.Monitor.MqttClientId)
	err = mqttClient.Connect()
	if err != nil {
		log.Fatal("MQTT broker not found... exiting.")
	}
	defer mqttClient.Disconnect()

	mqttClient.Subscribe("/heartbeat/#", monit.OnMessageReceivedHeartbeat)

	monit.RunMonitor(config.Monitor.ProcessNames)

	mqttClient.LoopForever()
}

//func sendStatusMessage(messageType mqtt2.StatusType, message ...any) {
//	messageTxt := fmt.Sprint(message...)
//
//	mqttMessage := &mqtt2.StatusPublishMessage{
//		Type:      messageType,
//		Status:    messageTxt,
//		Timestamp: time.Now(),
//	}
//
//	mqttClient.Publish(&mqtt2.Message{
//		Topic: "/status/monitor",
//		Value: mqttMessage,
//	})
//
//	if config.Monitor.LogStatusToConsole {
//		log.Println(messageType, ": ", messageTxt)
//	}
//}
