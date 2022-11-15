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

type VolumeChangePublishMessage struct {
	VolumeOffset float64 `json:"volumeOffset"`
}

type TrackChangePublishMessage struct {
	Direction int `json:"direction"`
}

type StatusPublishMessage struct {
	Type   StatusType `json:"type"`
	Status string     `json:"status"`
}

type PlayDonePublishMessage struct {
	Uid uint `json:"uid"`
}
