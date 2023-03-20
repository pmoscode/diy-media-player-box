package config

type Config struct {
	MqttBroker   MqttBroker   `yaml:"mqttBroker"`
	IoController IoController `yaml:"ioController"`
}

type MqttBroker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type IoController struct {
	MqttClientId       string  `yaml:"mqttClientId"`
	MockVolumeOffset   float64 `yaml:"mockVolumeOffset"`
	VolumeOffset       float64 `yaml:"volumeOffset"`
	PinVolumeUp        string  `yaml:"pinVolumeUp"`
	PinVolumeDown      string  `yaml:"pinVolumeDown"`
	PinTrackNext       string  `yaml:"pinTrackNext"`
	PinTrackPrev       string  `yaml:"pinTrackPrev"`
	LogStatusToConsole bool    `yaml:"logStatusToConsole"`
}
