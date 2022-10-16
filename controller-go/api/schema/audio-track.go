package schema

type AudioTrack struct {
	Track    uint   `json:"track"`
	Title    string `json:"title"`
	Length   string `json:"length"`
	FileName string `json:"filename"`
}
