package mqtt

import "time"

type StatusType string

const (
	Info  StatusType = "info"
	Error StatusType = "error"
	Warn  StatusType = "warn"
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
	Type      StatusType `json:"type" binding:"required"`
	Status    string     `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
}

type PlayDonePublishMessage struct {
	Uid uint `json:"uid"`
}
