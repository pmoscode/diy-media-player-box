#!/usr/bin/env bash

sshpass -p "raspberry" scp -r ./ pi@10.0.2.80:/home/pi/audio-player
sshpass -p "raspberry" ssh pi@10.0.2.80 "sudo apt install -y libasound2-dev"
sshpass -p "raspberry" ssh pi@10.0.2.80 "cd /home/pi/audio-player && /usr/local/go/bin/go mod tidy"
sshpass -p "raspberry" ssh pi@10.0.2.80 "cd /home/pi/audio-player && /usr/local/go/bin/go build -o audioPlayer"
