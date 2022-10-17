package api

type response struct {
	message string
}

type audioPlayerPublishMessage struct {
	Id        string   `json:"id"`
	TrackList []string `json:"trackList"`
}

type rfidReaderSubscribeMessage struct {
	CardId string `json:"cardId"`
}
