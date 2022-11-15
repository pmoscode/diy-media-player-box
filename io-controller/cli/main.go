package cli

import (
	"flag"
	"log"
)

func GetCliOptions() *Options {
	mqttBrokerIp := flag.String("mqtt-broker", "localhost", "Ip of MQTT broker")
	mqttClientId := flag.String("mqtt-client-id", "audio-player", "Client id for Mqtt connection")
	mockVolumeOffset := flag.Float64("mock-volume-offset", 1, "Volume offset to change (+/-)")
	logStatusToConsole := flag.Bool("log-console", false, "Log messages also to current std console")

	pinVolumeUp := flag.String("pin-volume-up", "12", "Pin on Raspi to control volume up changes")
	pinVolumeDown := flag.String("pin-volume-down", "13", "Pin on Raspi to control volume down changes")
	pinTrackNext := flag.String("pin-track-next", "16", "Pin on Raspi to control track next changes")
	pinTrackPrev := flag.String("pin-track-prev", "20", "Pin on Raspi to control track prev changes")
	flag.Parse()

	log.Println("Publishing / Subscribing to broker: ", *mqttBrokerIp)

	return &Options{
		MqttBrokerIp:       mqttBrokerIp,
		MqttClientId:       mqttClientId,
		MockVolumeOffset:   mockVolumeOffset,
		PinVolumeUp:        pinVolumeUp,
		PinVolumeDown:      pinVolumeDown,
		PinTrackNext:       pinTrackNext,
		PinTrackPrev:       pinTrackPrev,
		LogStatusToConsole: logStatusToConsole,
	}
}
