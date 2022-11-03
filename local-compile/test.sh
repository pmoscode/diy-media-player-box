#!/usr/bin/env bash
set -e


#docker run --rm -it --privileged \
#      -v $PWD:/go/src/diy-audio-player \
#      -w /go/src/diy-audio-player \
#      -v /var/run/docker.sock:/var/run/docker.sock \
#      -v $PWD/sysroot:/sysroot \
#      -e CGO_ENABLED=1 \
#      --entrypoint /bin/bash \
#      goreleaser/goreleaser-cross:v1.19.2

docker run --rm -it -v "$PWD":/usr/src/diy-media-box \
        -w /usr/src/diy-media-box/rfid-reader \
        diy-cross-compiler \
        /bin/bash -c 'go build -mod=readonly'
