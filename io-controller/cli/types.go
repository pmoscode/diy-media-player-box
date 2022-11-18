package cli

type Options struct {
	MqttBrokerIp       *string
	MqttClientId       *string
	MockVolumeOffset   *float64
	VolumeOffset       *float64
	PinVolumeUp        *string
	PinVolumeDown      *string
	PinTrackNext       *string
	PinTrackPrev       *string
	LogStatusToConsole *bool
}
