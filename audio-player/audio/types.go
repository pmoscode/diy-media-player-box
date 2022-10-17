package audio

type TracksSubscriptionMessage struct {
	Id        string   `json:"id" binding:"required"`
	TrackList []string `json:"trackList" binding:"required"`
}

type StatusPublishMessage struct {
	Status string `json:"status"`
}
