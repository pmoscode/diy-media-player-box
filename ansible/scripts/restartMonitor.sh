#!/usr/bin/env bash

PROCESS_NAME=monitor
COMMAND="./monitor >/dev/null 2>&1 &"

PID_PROCESS=$(pgrep -x $PROCESS_NAME)
echo "PID of '$PROCESS_NAME': $PID_PROCESS"

if [ -n "$PID_PROCESS" ]
then
  kill -9 "$PID_PROCESS"
fi

eval "$COMMAND"
