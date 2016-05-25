FROM golang:1.6.2-alpine

RUN go get  github.com/golang/lint/golint \
            github.com/tools/godep \
            github.com/laher/goxc

ENV USER root
WORKDIR /go/src/github.com/yuuki/ikachancat

ADD . /go/src/github.com/yuuki/ikachancat
