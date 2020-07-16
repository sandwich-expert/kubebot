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

# Ideally at this point we'd keep the image size down by doing a second "FROM ubuntu:20.04" and copy over kubebot and kubectl.
# However this results in segfaults in kubebot when connecting to Slack, so we just stick to the original golang image.
#FROM ubuntu:20.04
#COPY --from=0 /usr/local/bin/kubebot /usr/local/bin/kubebot
#COPY --from=0 /usr/local/bin/kubectl /usr/local/bin/kubectl

RUN ls -al /usr/local/bin \
  && kubectl version --client \
  && kubebot --help

CMD ["/usr/local/bin/kubebot"]
