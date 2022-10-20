package mqtt

type StatusType string

const (
	Info  StatusType = "info"
	Error StatusType = "error"
)

type Message struct {
	Topic string
	Value interface{}
}

type Response struct {
	Message string
}

type AudioPlayerPublishMessage struct {
	Id        uint     `json:"id" binding:"required"`
	TrackList []string `json:"trackList" binding:"required"`
}

type RfidReaderSubscribeMessage struct {
	CardId string `json:"cardId"`
}

type StatusPublishMessage struct {
	Type   StatusType `json:"type" binding:"required"`
	Status string     `json:"status"`
}

type PlayDoneSubscribeMessage struct {
	Uid uint `json:"uid"`
}
