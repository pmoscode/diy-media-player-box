package io

import (
	"github.com/pmoscode/go-common/mqtt"
	"io-controller/config"
	"log"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func NewOI(ioControllerConfig config.IoController, volumeChangeMessage func(volumeOffset float64), trackChangeMessage func(direction int), statusMessage func(messageType mqtt.StatusType, message ...any)) *IO {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	log.Println("Wiring pins to:")
	log.Println("\tVolDown:   ", "GPIO", ioControllerConfig.PinVolumeDown)
	log.Println("\tVolUp:     ", "GPIO", ioControllerConfig.PinVolumeUp)
	log.Println("\tTrackNext: ", "GPIO", ioControllerConfig.PinTrackNext)
	log.Println("\tTrackPrev: ", "GPIO", ioControllerConfig.PinTrackPrev)

	return &IO{
		sendStatusMessage:       statusMessage,
		sendVolumeChangeMessage: volumeChangeMessage,
		sendTrackChangeMessage:  trackChangeMessage,
		volDownPin:              gpioreg.ByName("GPIO" + ioControllerConfig.PinVolumeDown),
		volUpPin:                gpioreg.ByName("GPIO" + ioControllerConfig.PinVolumeUp),
		trackNextPin:            gpioreg.ByName("GPIO" + ioControllerConfig.PinTrackNext),
		trackPrevPin:            gpioreg.ByName("GPIO" + ioControllerConfig.PinTrackPrev),
		volumeOffset:            ioControllerConfig.VolumeOffset,
	}
}
