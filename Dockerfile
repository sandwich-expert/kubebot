FROM golang:1.14-alpine

# binutils: For 'strip'
RUN apk update \
  && apk add \
    binutils \
    curl \
    git \
  && rm -rf /var/cache/apk/*

# Download latest stable kubectl
RUN KUBECTL_VERSION=$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt) \
  && curl -L -o /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
  && strip /usr/local/bin/kubectl \
  && chmod +x /usr/local/bin/kubectl

# Get dependencies and build kubebot
ADD . /go/src/github.com/thelastpickle/kubebot
RUN go get github.com/thelastpickle/kubebot/ \
  && go build -o /usr/local/bin/kubebot github.com/thelastpickle/kubebot/ \
  && strip /usr/local/bin/kubebot

# Use alpine specifically to line up with golang-alpine above.
# Otherwise we see segfaults in kubebot when it's connecting to Slack.
FROM alpine
COPY --from=0 /usr/local/bin/kubebot /usr/local/bin/kubebot
COPY --from=0 /usr/local/bin/kubectl /usr/local/bin/kubectl

RUN ls -al /usr/local/bin \
  && kubectl version --client \
  && kubebot --help

CMD ["/usr/local/bin/kubebot"]
