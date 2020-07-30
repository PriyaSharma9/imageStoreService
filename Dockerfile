FROM golang:1.13-alpine as builder
ENV GOOS=linux
ENV GOARCH=amd64
ENV GIT_TERMINAL_PROMPT=1
# Argument 
ARG BRANCH=master
# steps for git clone
RUN git clone -b $BRANCH https://github.com/PriyaSharma9/imageStoreService
WORKDIR imageStoreService

# Dependencies need to be added in the container

RUN go build -o imageStoreService main.go
# Copying Local files to the container
# Monitoring script
FROM alpine:3.11.3
MAINTAINER priya.sharma6693@gmail.com
RUN mkdir -p /usr/share/imageStoreService/

COPY isDocker.sh monitor_imageStoreService.sh /
COPY --from=builder /go/imageStoreService/imageStoreService /usr/share/imageStoreService/imageStoreService

# Packages
#RUN apk cache clean && \
RUN apk update && \
    apk add -f libcurl  && \
    apk add -f jsoncpp && \
    apk add -f libevent-dev && \
    apk add -f curl && \
    apk add -f net-tools && \
    apk add -f iputils && \
    apk add -f iproute2 && \
    apk add -f vim && \
    apk add -f ndisc6 && \
    apk add -f busybox-extras && \
    apk add -f netcat-openbsd && \
    apk add -f apk-cron && \
    apk add -f openssl && \
    apk add  bash && \
    chmod 777 /isDocker.sh  && \
    chmod 777 /monitor_imageStoreService.sh
ENTRYPOINT ["/isDocker.sh"]
