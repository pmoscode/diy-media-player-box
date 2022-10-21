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

type TracksSubscriptionMessage struct {
	Id        uint     `json:"id" binding:"required"`
	TrackList []string `json:"trackList" binding:"required"`
}

type VolumeChangeSubscriptionMessage struct {
	VolumeOffset float64 `json:"volumeOffset"`
}

type StatusPublishMessage struct {
	Type   StatusType `json:"type"`
	Status string     `json:"status"`
}

type PlayDonePublishMessage struct {
	Uid uint `json:"uid"`
}
