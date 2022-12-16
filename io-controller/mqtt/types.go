package mqtt

type VolumeChangePublishMessage struct {
	VolumeOffset float64 `json:"volumeOffset"`
}

type TrackChangePublishMessage struct {
	Direction int `json:"direction"`
}
