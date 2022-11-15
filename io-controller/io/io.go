package io

import (
	"fmt"
	"io-controller/cli"
	"io-controller/mqtt"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"time"
)

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
	log.Println(gpioreg.All())

	if err := i.volUpPin.In(gpio.Float, gpio.NoEdge); err != nil {
		log.Println("error occurred")
		log.Fatal(err)
	}

	fmt.Printf("%s: %s\n", i.volUpPin, i.volUpPin.Function())

	i.sendStatusMessage(mqtt.Info, "Pin configured as input...")

	for {
		i.sendStatusMessage(mqtt.Info, "Waiting for input...")
		//read := i.volUpPin.WaitForEdge(1 * time.Second)
		//if read {
		time.Sleep(1 * time.Second)
		i.sendStatusMessage(mqtt.Info, "... Pin is: ", fmt.Sprintf("-> %s\n", i.volUpPin.Read().String()))
		//} else {
		//	i.sendStatusMessage(mqtt.Info, "... Pin ", i.volUpPin.Name(), " was not triggered...")
		//}
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
