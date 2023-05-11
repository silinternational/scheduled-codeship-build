FROM golang:1.17

RUN curl -o- -L https://slss.io/install | VERSION=3.3.0 bash

# Install 1Password CLI
RUN apt-get update && apt-get install -y unzip
RUN curl -sSfo op.zip \
  https://cache.agilebits.com/dist/1P/op2/pkg/v2.17.0/op_linux_amd64_v2.17.0.zip \
  && unzip -od /usr/local/bin/ op.zip \
  && rm op.zip

RUN alias ll="ls -al"

# Copy in source and install deps
WORKDIR /src
COPY ./ /src/
RUN go get ./...
