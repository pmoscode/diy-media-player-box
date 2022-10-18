package audio

type TracksSubscriptionMessage struct {
	Id        uint     `json:"id" binding:"required"`
	TrackList []string `json:"trackList" binding:"required"`
}

type StatusPublishMessage struct {
	Status string `json:"status"`
}

type PlayDonePublishMessage struct {
	Uid uint `json:"uid"`
}
