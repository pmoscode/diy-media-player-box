#!/usr/bin/env python3

import os
import time

import RPi.GPIO as GPIO
import mfrc522
import requests

reader = mfrc522.MFRC522()

last_id = None
last_status = -1
remove_counter = 0
remove_ok_threshold = 2

controller_port = os.getenv("CONTROLLER_PORT", 2020)


def start():
    global last_status, remove_counter, last_id

    try:
        while True:
            # Scan for cards
            (status, TagType) = reader.MFRC522_Request(reader.PICC_REQIDL)

            if last_status == -1 or last_status != status:
                last_status = status

            # If a card is found
            if status == reader.MI_OK:
                # Get the UID of the card
                (status, uid) = reader.MFRC522_Anticoll()
                string_uid = [str(num) for num in uid]
                tag_id = "".join(string_uid)
                if tag_id and tag_id != last_id:
                    print("New card present...", tag_id)
                    send_play_request(tag_id)
                    last_id = tag_id
                remove_counter = 0
            if status == reader.MI_ERR and last_id is not None:
                remove_counter += 1
                if remove_counter >= remove_ok_threshold:
                    print("Card removed...")
                    send_pause_request()
                    last_id = None
                    remove_counter = 0
                    last_status = -1
            time.sleep(0.5)
    finally:
        GPIO.cleanup()


def send_play_request(tag_id: str):
    global controller_port
    api_endpoint = f"http://localhost:{controller_port}/audio-books/{tag_id}/play"

    send_request(api_endpoint)


def send_pause_request():
    global controller_port
    api_endpoint = f"http://localhost:{controller_port}/audio-books/pause"

    send_request(api_endpoint)


def send_request(api_endpoint: str):
    r = requests.post(url=api_endpoint)

    response = r.text
    print("Response: %s" % response)


if __name__ == "__main__":
    start()
