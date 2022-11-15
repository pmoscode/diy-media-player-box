package io

import (
	"io-controller/cli"
	"io-controller/mqtt"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"time"
)

const volumeOffset = 0.5

type IO struct {
	sendVolumeChangeMessage func(volumeOffset float64)
	sendTrackChangeMessage  func(direction int)
	sendStatusMessage       func(messageType mqtt.StatusType, message ...any)
	volUpPin                gpio.PinIn
	volDownPin              gpio.PinIn
	trackNext               gpio.PinIn
	trackPrev               gpio.PinIn
}

func (i *IO) Run() {
	i.sendStatusMessage(mqtt.Info, "Configuring pins...")
	i.setupPins()
	i.sendStatusMessage(mqtt.Info, "...Pins configured")

	i.sendStatusMessage(mqtt.Info, "Starting loop...")
	for {
		i.checkVolumeStates()
		i.checkTrackStates()

		time.Sleep(1 * time.Second)
	}
}

func (i *IO) setupPins() {
	if err := i.volUpPin.In(gpio.Float, gpio.NoEdge); err != nil {
		log.Fatal(err)
	}
	if err := i.volDownPin.In(gpio.Float, gpio.NoEdge); err != nil {
		log.Fatal(err)
	}
	if err := i.trackNext.In(gpio.Float, gpio.NoEdge); err != nil {
		log.Fatal(err)
	}
	if err := i.trackPrev.In(gpio.Float, gpio.NoEdge); err != nil {
		log.Fatal(err)
	}
}

func (i *IO) checkVolumeStates() {
	if i.volUpPin.Read() == gpio.High {
		i.sendStatusMessage(mqtt.Info, "Volume up button pressed")
		i.sendVolumeChangeMessage(volumeOffset)
	}

	if i.volDownPin.Read() == gpio.High {
		i.sendStatusMessage(mqtt.Info, "Volume down button pressed")
		i.sendVolumeChangeMessage(-volumeOffset)
	}
}

func (i *IO) checkTrackStates() {
	if i.trackNext.Read() == gpio.High {
		i.sendStatusMessage(mqtt.Info, "Next track button pressed")
		i.sendTrackChangeMessage(1)
	}

	if i.trackPrev.Read() == gpio.High {
		i.sendStatusMessage(mqtt.Info, "Previous track button pressed")
		i.sendTrackChangeMessage(-1)
	}
}

func NewOI(cliOptions *cli.Options, volumeChangeMessage func(volumeOffset float64), trackChangeMessage func(direction int), statusMessage func(messageType mqtt.StatusType, message ...any)) *IO {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	return &IO{
		sendStatusMessage:       statusMessage,
		sendVolumeChangeMessage: volumeChangeMessage,
		sendTrackChangeMessage:  trackChangeMessage,
		volDownPin:              gpioreg.ByName("GPIO" + *cliOptions.PinVolumeDown),
		volUpPin:                gpioreg.ByName("GPIO" + *cliOptions.PinVolumeUp),
		trackNext:               gpioreg.ByName("GPIO" + *cliOptions.PinTrackNext),
		trackPrev:               gpioreg.ByName("GPIO" + *cliOptions.PinTrackPrev),
	}
}
