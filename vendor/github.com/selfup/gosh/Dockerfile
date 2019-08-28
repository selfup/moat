FROM golang:alpine

ENV GOSH src/github.com/selfup/gosh

RUN mkdir -p go/src/github.com/selfup/gosh

COPY . $GOPATH/$GOSH

WORKDIR $GOPATH/$GOSH

RUN go test
