package rfid

import (
	"gitlab.com/pmoscodegrp/common/mqtt"
)

type Mock struct {
	cardId            string
	sendCardIdMessage func(cardId string)
	sendStatusMessage func(messageType mqtt.StatusType, message ...any)
}

func (m *Mock) Run() {
	m.sendCardIdMessage(m.cardId)
	m.sendStatusMessage(mqtt.Info, "Mock CardId '", m.cardId, "' send to controller...")
}

func NewMock(cardId string, cardIdMessage func(cardId string), statusMessage func(messageType mqtt.StatusType, message ...any)) *Mock {
	return &Mock{cardId: cardId, sendStatusMessage: statusMessage, sendCardIdMessage: cardIdMessage}
}
