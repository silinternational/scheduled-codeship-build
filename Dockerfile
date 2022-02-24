FROM golang:1.17

RUN curl -o- -L https://slss.io/install | VERSION=3.3.0 bash

RUN alias ll="ls -al"

# Copy in source and install deps
WORKDIR /src
COPY ./ /src/
RUN go get ./...
