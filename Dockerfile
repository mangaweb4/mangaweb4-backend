#Stage 1 -- building executable
FROM golang:1.24-alpine AS builder1

WORKDIR /go/src/mangaweb
COPY . .

ARG VERSION=Development
RUN apk add git
RUN go get -d -v ./...
RUN go build -v -ldflags="-X 'main.versionString=$VERSION' " -o mangaweb .

# Stage 2 -- build the target image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder1 /go/src/mangaweb/mangaweb ./

CMD ["./mangaweb"]