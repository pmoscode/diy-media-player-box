#!/usr/bin/env bash

docker run -it --rm --name diy_media_box_controller -v "$PWD":/usr/src/app -w /usr/src/app -p 2020:2020 node:14.19-alpine3.15 node server.js
