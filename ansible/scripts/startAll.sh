#!/usr/bin/env bash

echo "Starting rfid-reader..."
./rfid-reader >/dev/null 2>&1 &
sleep 2

echo "Starting audio-player..."
./audio-player --buffer-sample-rate 700 >/dev/null 2>&1 &
sleep 2

echo "Starting controller..."
./controller >/dev/null 2>&1 &
sleep 4

echo "Starting io-controller..."
./io-controller >/dev/null 2>&1 &
sleep 2

echo "Starting logger..."
./logger >/dev/null 2>&1 &
