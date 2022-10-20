package mqtt

type StatusType string

const (
	Info  StatusType = "info"
	Error StatusType = "error"
)

type CardIdPublishMessage struct {
	CardId string `json:"cardId"`
}

type StatusPublishMessage struct {
	Type   StatusType `json:"type"`
	Status string     `json:"status"`
}

type Message struct {
	Topic string
	Value interface{}
}
