#!/usr/bin/env bash

echo "Starting logger..."
./logger >/dev/null 2>&1 &
sleep 2

echo "Starting rfid-reader..."
./rfid-reader --remove-threshold 2 >/dev/null 2>&1 &
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

echo "Starting monitor..."
./monitor >/dev/null 2>&1 &
