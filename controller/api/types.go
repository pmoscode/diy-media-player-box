package api

type response struct {
	message string
}

type audioPlayerPublishMessage struct {
	Id        uint     `json:"id"`
	TrackList []string `json:"trackList"`
}

type rfidReaderSubscribeMessage struct {
	CardId string `json:"cardId"`
}

type statusPublishMessage struct {
	Status string `json:"status"`
}

type PlayDoneSubscribeMessage struct {
	Uid uint `json:"uid"`
}
