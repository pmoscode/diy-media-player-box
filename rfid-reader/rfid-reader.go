package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

var (
	rfid           *mfrc522.Dev
	lastId         = ""
	removeCounter  = 0
	controllerPort = getEnv("CONTROLLER_PORT", "2020")
)

const removeOkThreshold = 2

func main() {
	var err error
	err = nil

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	p, errOpen := spireg.Open("")
	if errOpen != nil {
		log.Fatal(err)
	}
	defer p.Close()

	rfid, err = mfrc522.NewSPI(p, rpi.P1_22, rpi.P1_18)
	if err != nil {
		log.Fatal(err)
	}

	defer rfid.Halt()

	rfid.SetAntennaGain(6)

	for true {
		search()
	}

}

func search() {
	timedOut := false
	cb := make(chan []byte)
	timer := time.NewTimer(time.Second)

	defer func() {
		timer.Stop()
		timedOut = true
		close(cb)
	}()

	go func() {
		for {
			uid, err := rfid.ReadUID(time.Second)

			if timedOut {
				return
			}

			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			cb <- uid
			return
		}
	}()

	for {
		// fmt.Print("\033[G")
		select {
		case <-timer.C:
			if lastId != "" {
				removeCounter++
				if removeCounter >= removeOkThreshold {
					log.Print("Card removed...")
					sendRequest(fmt.Sprintf("http://localhost:%s/api/audio-books/pause", controllerPort))
					lastId = ""
					removeCounter = 0
				}
			}

			// fmt.Print("\033[A")
			return
		case data := <-cb:
			cardId := hex.EncodeToString(data)
			if cardId != lastId {
				fmt.Println("New card present: ", cardId)
				lastId = cardId
				sendRequest(fmt.Sprintf("http://localhost:%s/api/audio-books/%s/play", controllerPort, cardId))
			}
			removeCounter = 0
			time.Sleep(time.Second)
			// fmt.Print("\033[A")
			return
		}
	}
}

func getEnv(name string, defaultValue string) string {
	env, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return env
}

func sendRequest(apiEndpoint string) {
	fmt.Println("Sending request to: ", apiEndpoint)

	resp, err := http.Post(apiEndpoint, "application/json", nil)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println("Code: ", resp.StatusCode, " ## Response: ", res["json"])
}
