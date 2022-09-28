#!/usr/bin/env bash

sshpass -p "raspberry" scp -r startAll.sh pi@10.0.2.80:/home/pi
