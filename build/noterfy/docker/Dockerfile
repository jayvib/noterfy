# Copyright 2021 Jayson Vibandor. All right reserved.

# STEP 1: Build the executable
FROM golang:1.16.3-alpine3.13 as builder
RUN apk --update add git make
WORKDIR /go/src/noterfy
COPY . .
RUN make mod
RUN make build-noterfy

# STEP 2: Distribute the executable
FROM alpine:3.7
RUN apk add --no-cache \
			bash \
			ca-certificates \
			curl \
ARG NOTERFY_BUILD_COMMIT
ARG NOTERFY_VERSION
LABEL "com.noterfy.description"="The server engine of the Noterfy Application."
LABEL "com.noterfy.maintainer"="Jayson Vibandor jayson.vibandor@gmail.com"
LABEL "com.noterfy.version"=$NOTERFY_VERSION
LABEL "com.noterfy.build_commit"=$NOTERFY_BUILD_COMMIT
RUN mkdir /home/noterfy
EXPOSE 50001/tcp
COPY --from=builder /go/src/noterfy/bin/noterfy.linux /home/noterfy/
RUN chmod +x /home/noterfy/noterfy.linux
CMD [ "/home/noterfy/noterfy.linux" ]
