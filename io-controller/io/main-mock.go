package io

import "io-controller/mqtt"

type Mock struct {
	volumeOffset            float64
	sendVolumeChangeMessage func(volumeOffset float64)
	sendStatusMessage       func(messageType mqtt.StatusType, message ...any)
}

func (m *Mock) Run() {
	m.sendVolumeChangeMessage(m.volumeOffset)
	m.sendStatusMessage(mqtt.Info, "Mock VolumeOffset '", m.volumeOffset, "' send to controller...")
}

func NewMockOI(volumeOffset float64, volumeChangeMessage func(volumeOffset float64), statusMessage func(messageType mqtt.StatusType, message ...any)) *Mock {
	return &Mock{volumeOffset: volumeOffset, sendStatusMessage: statusMessage, sendVolumeChangeMessage: volumeChangeMessage}
}
