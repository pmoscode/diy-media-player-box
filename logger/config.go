package main

type Config struct {
	MqttBroker ConfigMqttBroker `yaml:"mqttBroker"`
	Logger     ConfigLogger     `yaml:"logger"`
}

type ConfigMqttBroker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigLogger struct {
	MqttClientId                string `yaml:"mqttClientId"`
	LogRotationPeriodAfterBytes int    `yaml:"logRotationPeriodAfterBytes"`
	LogStatusToConsole          bool   `yaml:"logStatusToConsole"`
	FileName                    string `yaml:"fileName"`
	MqttSubscriptionTopic       string `yaml:"mqttSubscriptionTopic"`
}
