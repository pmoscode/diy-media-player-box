#!/usr/bin/env bash

function check() {
  PROCESS_NAME=$1
  COMMAND=$2
  PID_PROCESS=$(pgrep $PROCESS_NAME)
  echo "PID of '$PROCESS_NAME': $PID_PROCESS"
  [  -z "$PID_PROCESS" ] && echo "Starting process: $COMMAND"; eval "$COMMAND" || echo "$PROCESS_NAME not already running"
}

echo "Checking rfid-reader..."
check rfid-reader "./rfid-reader >/dev/null 2>&1 &"

echo "Checking audio-player..."
check audio-player "./audio-player --buffer-sample-rate 700 >/dev/null 2>&1 &"

echo "Checking controller..."
check controller "./controller >/dev/null 2>&1 &"

echo "Checking io-controller..."
check io-controller "./io-controller >/dev/null 2>&1 &"
