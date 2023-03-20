package main

type Config struct {
	MqttBroker ConfigMqttBroker `yaml:"mqttBroker"`
	Monitor    ConfigMonitor    `yaml:"monitor"`
}

type ConfigMqttBroker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigMonitor struct {
	MqttClientId       string `yaml:"mqttClientId"`
	ProcessNames       string `yaml:"processNames"`
	LogStatusToConsole bool   `yaml:"logStatusToConsole"`
}
