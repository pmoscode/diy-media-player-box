#!/usr/bin/env bash

sshpass -p "raspberry" scp -r ./ pi@10.0.2.80:/home/pi/rfid-reader
# sshpass -p "raspberry" ssh pi@10.0.2.80 "cd /home/pi/rfid-reader && /usr/local/go/bin/go mod tidy"
sshpass -p "raspberry" ssh pi@10.0.2.80 "cd /home/pi/rfid-reader && /usr/local/go/bin/go build -o rfid-reader"
