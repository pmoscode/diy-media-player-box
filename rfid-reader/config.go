package main

type Config struct {
	MqttBroker ConfigMqttBroker `yaml:"mqttBroker"`
	RfidReader ConfigRfidReader `yaml:"rfidReader"`
}

type ConfigMqttBroker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigRfidReader struct {
	MqttClientId       string `yaml:"mqttClientId"`
	MockCardId         string `yaml:"mockCardId"`
	LogStatusToConsole bool   `yaml:"logStatusToConsole"`
	RemoveThreshold    int    `yaml:"removeThreshold"`
}
