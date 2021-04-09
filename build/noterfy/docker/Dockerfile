# Copyright 2021 Jayson Vibandor. All right reserved.

# STEP 1: Build the executable
FROM golang:1.16.3-alpine3.13 as builder
MAINTAINER "Jayson Vibandor" "jayson.vibandor@gmail.com"
RUN apk update && \
		apk upgrade && \
	  apk --update add git gcc make
WORKDIR /go/src/noterfy
COPY . .
RUN make mod
RUN make build-noterfy

# STEP 2: Distribute the executable
FROM alpine:3.7
RUN apk update && \
		apk upgrade && \
		apk --update --no-cache add tzdata && \
		apk add --no-cache bash ca-certificates openssh
RUN set -ex && apk add --no-cache --virtual bash musl-dev openssl
RUN mkdir /home/noterfy
EXPOSE 50001
COPY --from=builder /go/src/noterfy/bin/noterfy.linux /home/noterfy/
RUN chmod +x /home/noterfy/noterfy.linux
CMD [ "/home/noterfy/noterfy.linux" ]
