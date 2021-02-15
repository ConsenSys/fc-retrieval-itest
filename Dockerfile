# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.15-alpine

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/src/github.com/ConsenSys/fc-retrieval-itest
COPY . .
# Remove any cached dependancies. TODO is this really needed?
RUN go clean -modcache

CMD go test -v ./...