package mqtt

import "time"

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
	Type      StatusType `json:"type" binding:"required"`
	Status    string     `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
}
