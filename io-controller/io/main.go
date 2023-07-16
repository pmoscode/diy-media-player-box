package io

import (
	"github.com/pmoscode/go-common/mqtt"
	"io-controller/config"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"time"
)

const waitForEdgeTime = 100 // In milliseconds

type IO struct {
	sendVolumeChangeMessage func(volumeOffset float64)
	sendTrackChangeMessage  func(direction int)
	sendStatusMessage       func(messageType mqtt.StatusType, message ...any)
	volUpPin                gpio.PinIn
	volUpLastState          gpio.Level
	volDownPin              gpio.PinIn
	volDownLastState        gpio.Level
	trackNextPin            gpio.PinIn
	trackNextLastState      gpio.Level
	trackPrevPin            gpio.PinIn
	trackPrevLastState      gpio.Level
	volumeOffset            float64
}

func (i *IO) Run() {
	i.sendStatusMessage(mqtt.Info, "Configuring pins and states...")
	i.setupPins()
	i.setupStates()
	i.sendStatusMessage(mqtt.Info, "...done")

	i.sendStatusMessage(mqtt.Info, "Starting loop...")
	for {
		i.checkVolumeStates()
		i.checkTrackStates()
	}
}

func (i *IO) setupStates() {
	i.volUpLastState = gpio.Low
	i.volDownLastState = gpio.Low
	i.trackNextLastState = gpio.Low
	i.trackPrevLastState = gpio.Low
}

func (i *IO) setupPins() {
	if err := i.volUpPin.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}
	if err := i.volDownPin.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}
	if err := i.trackNextPin.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}
	if err := i.trackPrevPin.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}
}

func (i *IO) checkVolumeStates() {
	up := i.volUpPin.WaitForEdge(waitForEdgeTime * time.Millisecond)
	if up && i.volUpLastState == gpio.Low {
		i.volUpLastState = gpio.High
		i.sendStatusMessage(mqtt.Info, "Volume up button pressed")
		i.sendVolumeChangeMessage(i.volumeOffset)
	} else if !up {
		i.volUpLastState = gpio.Low
	}

	down := i.volDownPin.WaitForEdge(waitForEdgeTime * time.Millisecond)
	if down && i.volDownLastState == gpio.Low {
		i.volDownLastState = gpio.High
		i.sendStatusMessage(mqtt.Info, "Volume down button pressed")
		i.sendVolumeChangeMessage(-i.volumeOffset)
	} else if !up {
		i.volDownLastState = gpio.Low
	}
}

func (i *IO) checkTrackStates() {
	next := i.trackNextPin.WaitForEdge(waitForEdgeTime * time.Millisecond)
	if next && i.trackNextLastState == gpio.Low {
		i.trackNextLastState = gpio.High
		i.sendStatusMessage(mqtt.Info, "Next track button pressed")
		i.sendTrackChangeMessage(1)
	} else if !next {
		i.trackNextLastState = gpio.Low
	}

	prev := i.trackPrevPin.WaitForEdge(waitForEdgeTime * time.Millisecond)
	if prev && i.trackPrevLastState == gpio.Low {
		i.trackPrevLastState = gpio.High
		i.sendStatusMessage(mqtt.Info, "Previous track button pressed")
		i.sendTrackChangeMessage(-1)
	} else if !next {
		i.trackPrevLastState = gpio.Low
	}
}

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
