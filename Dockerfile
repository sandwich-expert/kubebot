# upstream go project changed and build no longer works, as workaround, pulling previous image that already include dependencies
FROM thelastpickle/kubebot:0bd34d7

RUN set -e && \
    curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl && \
    kubectl version --client

RUN rm -rf /go/src/app/

ADD . /go/src/app/

RUN go build -o app *.go

CMD ["/go/src/app/app"]
