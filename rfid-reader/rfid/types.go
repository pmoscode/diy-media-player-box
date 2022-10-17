package rfid

type CardIdPublishMessage struct {
	CardId string `json:"cardId"`
}

type StatusPublishMessage struct {
	Status string `json:"status"`
}
