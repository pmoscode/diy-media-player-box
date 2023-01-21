#!/usr/bin/env bash

function check() {
  PROCESS_NAME=$1
  COMMAND=$2
  PID_PROCESS=$(pgrep -x $PROCESS_NAME)
  echo "PID of '$PROCESS_NAME': $PID_PROCESS"

  if [  -z "$PID_PROCESS" ]
  then
    echo "Starting process: $COMMAND"
    eval "$COMMAND"
  else
    echo "$PROCESS_NAME already running"
  fi
}

echo "Checking logger..."
check logger "./logger >/dev/null 2>&1 &"

echo "Checking rfid-reader..."
check rfid-reader "./rfid-reader --remove-threshold 2 >/dev/null 2>&1 &"

echo "Checking audio-player..."
check audio-player "./audio-player --buffer-sample-rate 700 >/dev/null 2>&1 &"

echo "Checking controller..."
check controller "./controller >/dev/null 2>&1 &"

echo "Checking io-controller..."
check io-controller "./io-controller >/dev/null 2>&1 &"

echo "Checking monitor..."
check monitor "./monitor >/dev/null 2>&1 &"
