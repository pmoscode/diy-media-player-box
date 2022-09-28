#!/usr/bin/env bash
set -e

echo "Setting environment"
export PORT=2020
export AUDIO_PLAYER_SAMPLE_RATE_FACTOR=1
export AUDIO_PLAYER_BIND_ADDRESS=localhost:8080
export GIN_MODE=release

echo "Start Rfid-Reader"
./rfid-reader/rfid-reader &

echo "Start Audio-Player"
./audio-player/audioPlayer &

echo "Start Controller"
cd controller/
/usr/local/bin/node server.js &
cd ..
