package audio

type audioRequestInput struct {
	Uid       string   `json:"uid" binding:"required"`
	TrackList []string `json:"trackList" binding:"required"`
}
