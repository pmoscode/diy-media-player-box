#!/usr/bin/env bash

sshpass -p "raspberry" scp -r app/ ui/ .env package.json server.js pi@10.0.2.80:/home/pi/controller
#sshpass -p "raspberry" ssh pi@10.0.2.80 "cd /home/pi/controller && mkdir database"
