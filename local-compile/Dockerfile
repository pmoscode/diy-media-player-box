FROM golang:1.23.0-bookworm
LABEL os=linux
LABEL arch=armhf

ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=6
ENV CGO_ENABLED=1
ENV CC=arm-linux-gnueabihf-gcc
ENV CXX=arm-linux-gnueabihf-g++
ENV PATH="/go/bin/${GOOS}_${GOARCH}:${PATH}"
ENV PKG_CONFIG_PATH=/usr/lib/arm-linux-gnueabihf/pkgconfig

# install build & runtime dependencies
RUN dpkg --add-architecture armhf \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        protobuf-compiler \
        upx \
        gcc-arm-linux-gnueabihf \
        libc6-dev-armhf-cross \
        pkg-config \
        alsa-utils:armhf \
        libasound2-dev:armhf \
    && rm -rf /var/lib/apt/lists/*

## install build dependencies (code generators)
#RUN GOARCH=amd64 go install github.com/gogo/protobuf/protoc-gen-gofast@latest \
#    && GOARCH=amd64 go install github.com/GeertJohan/go.rice/rice@latest \
#    && GOARCH=amd64 go install github.com/micro/protoc-gen-micro@latest
