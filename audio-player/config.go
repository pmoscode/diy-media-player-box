package main

type Config struct {
	MqttBroker  ConfigMqttBroker  `yaml:"mqttBroker"`
	AudioPlayer ConfigAudioPlayer `yaml:"audioPlayer"`
}

type ConfigMqttBroker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigAudioPlayer struct {
	MqttClientId       string `yaml:"mqttClientId"`
	SampleRateFactor   int    `yaml:"sampleRateFactor"`
	LogStatusToConsole bool   `yaml:"logStatusToConsole"`
}
