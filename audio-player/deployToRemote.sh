#!/usr/bin/env bash

export GOOS=linux
export GOARCH=arm
export GOARM=5

go build -o testRpi
sshpass -p "raspberry" scp testRpi pi@10.0.2.80:/home/pi/audio-player
