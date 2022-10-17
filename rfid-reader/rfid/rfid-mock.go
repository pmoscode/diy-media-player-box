package rfid

import (
	"log"
	"rfid-reader/mqtt"
)

type Mock struct {
	cardId     string
	mqttClient *mqtt.Client
}

func (m *Mock) Run() {
	log.Println("Sending mock message...")
	message := &CardIdPublishMessage{
		CardId: m.cardId,
	}

	m.mqttClient.SendMessage(&mqtt.Message{
		Topic: "/rfidReader/cardId",
		Value: message,
	})
}

func NewMock(cardId *string, mqttClient *mqtt.Client) *Mock {
	return &Mock{*cardId, mqttClient}
}
