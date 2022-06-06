#!/usr/bin/env bash

sshpass -p "raspberry" scp rfid-reader.py pi@10.0.2.80:/home/pi
sshpass -p "raspberry" scp requirements.txt pi@10.0.2.80:/home/pi
