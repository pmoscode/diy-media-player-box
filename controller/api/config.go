package api

type Config struct {
	MqttBroker ConfigMqttBroker `yaml:"mqttBroker"`
	Controller ConfigController `yaml:"controller"`
}

type ConfigMqttBroker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigController struct {
	MqttClientId string `yaml:"mqttClientId"`
}
