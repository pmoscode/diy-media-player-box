#!/usr/bin/env python3

import time

import RPi.GPIO as GPIO
import mfrc522
import requests

reader = mfrc522.MFRC522()

last_id = None
last_status = -1
remove_counter = 0
remove_ok_threshold = 2


def start():
    global last_status, remove_counter, last_id

    try:
        while True:
            # Scan for cards
            (status, TagType) = reader.MFRC522_Request(reader.PICC_REQIDL)
            # print("Reading tag: ", status, " # Tag Type: ", TagType)

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
                    send_tag_id(tag_id)
                    last_id = tag_id
                remove_counter = 0
            if status == reader.MI_ERR and last_id is not None:
                remove_counter += 1
                if remove_counter >= remove_ok_threshold:
                    print("Card removed...")
                    send_tag_id("")
                    last_id = None
                    remove_counter = 0
                    last_status = -1
            time.sleep(0.5)
    finally:
        GPIO.cleanup()


def send_tag_id(tag_id: str):
    api_endpoint = f"http://localhost:2020/audio-books/{tag_id}/play"

    r = requests.post(url=api_endpoint)

    response = r.text
    print("Response: %s" % response)


if __name__ == "__main__":
    start()
