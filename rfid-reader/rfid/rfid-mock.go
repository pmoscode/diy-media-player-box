package rfid

import (
	"log"
	"rfid-reader/mqtt"
)

type Mock struct {
	mqttClient *mqtt.Client
}

func (m *Mock) Run() {
	// To be implemented
	log.Println("Mock not implemented yet...")
	//m.mqttClient.LoopForever()
}

func NewMock(mqttClient *mqtt.Client) *Mock {
	return &Mock{mqttClient}
}
