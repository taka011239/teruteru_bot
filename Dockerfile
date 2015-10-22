FROM golang

ADD . /go/src/github.com/taka011239/teruteru_bot

RUN GO15VENDOREXPERIMENT=1 go install github.com/taka011239/teruteru_bot

COPY config.tml /go/

ENTRYPOINT /go/bin/teruteru_bot
