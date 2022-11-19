#!/usr/bin/env bash

function terminate() {
  PID=$(pgrep $1)
  echo "PID of '$1': $PID"
  [  -z "$PID" ] && echo "$1 not running" || kill -15 "$PID"
}

terminate audio-player
terminate rfid-reader
terminate io-controller
terminate controller
