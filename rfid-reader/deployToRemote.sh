#!/usr/bin/env bash

sshpass -p "raspberry" scp rfid-reader.py requirements.txt pi@10.0.2.80:/home/pi/rfid-reader
