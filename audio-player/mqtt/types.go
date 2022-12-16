package mqtt

type TracksSubscriptionMessage struct {
	Id        uint     `json:"id" binding:"required"`
	TrackList []string `json:"trackList" binding:"required"`
}

type VolumeChangeSubscriptionMessage struct {
	VolumeOffset float64 `json:"volumeOffset"`
}

type PlayDonePublishMessage struct {
	Uid uint `json:"uid"`
}
