package mqtt

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

type PlayDoneSubscribeMessage struct {
	Uid uint `json:"uid"`
}
