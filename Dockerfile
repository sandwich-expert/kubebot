# upstream go project changed and build no longer works, as workaround, pulling previous image that already include dependencies
FROM thelastpickle/kubebot:0bd34d7

RUN rm -rf /go/src/app/

ADD . /go/src/app/

RUN go build -o app *.go

CMD ["/go/src/app/app"]
