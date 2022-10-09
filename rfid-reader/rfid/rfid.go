package rfid

import (
	"encoding/hex"
	"log"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"strings"
	"time"
)

const removeOkThreshold = 2

type Rfid struct {
	rfid              *mfrc522.Dev
	lastId            string
	removeCounter     int
	sendStatusMessage func(message string)
	sendCardIdMessage func(cardId string)
}

func (r *Rfid) Run() {
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

	r.rfid, err = mfrc522.NewSPI(p, rpi.P1_22, rpi.P1_18)
	if err != nil {
		log.Fatal(err)
	}

	defer r.rfid.Halt()

	r.rfid.SetAntennaGain(6)

	cb := make(chan []byte)
	defer func() {
		close(cb)
	}()

	r.listen(cb)
}

func (r *Rfid) listen(cb chan []byte) {
	timedOut := false

	defer func() {
		timedOut = true
	}()

	go func() {
		for {
			uid, err := r.rfid.ReadUID(time.Second)

			if timedOut {
				return
			}

			if err != nil {
				if !strings.Contains(err.Error(), "timeout waiting for IRQ edge") {
					continue
				}
			}

			cb <- uid
			if len(uid) > 0 {
				time.Sleep(time.Second)
			}
		}
	}()

	for {
		select {
		case data := <-cb:
			cardId := hex.EncodeToString(data)
			//log.Println("Card found: ", cardId)
			if cardId != "" {
				if cardId != r.lastId {
					log.Println("New card present: ", cardId)
					r.sendCardIdMessage(cardId)
					r.sendStatusMessage("New card present: " + cardId)
					r.lastId = cardId
				}
				r.removeCounter = 0
			} else {
				if r.lastId != "" {
					r.removeCounter++
					if r.removeCounter >= removeOkThreshold {
						log.Print("Card removed...")
						r.sendCardIdMessage("")
						r.sendStatusMessage("Card removed: " + r.lastId)
						r.lastId = ""
						r.removeCounter = 0
					}
				}
			}
		}
	}
}

func NewRfid(cardIdMessage func(message string), statusMessage func(cardId string)) *Rfid {
	return &Rfid{sendStatusMessage: statusMessage, sendCardIdMessage: cardIdMessage}
}
