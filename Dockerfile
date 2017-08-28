FROM golang:1.8.3

RUN wget http://storage.googleapis.com/kubernetes-release/release/v1.7.3/bin/linux/amd64/kubectl -O /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app/

RUN set -x && \
    go get github.com/nlopes/slack && \
    go get github.com/go-chat-bot/bot
RUN go build -o app *.go

CMD ["app"]
