#Stage 1 -- building executable
FROM golang:1.24-alpine AS builder1

WORKDIR /go/src/mangaweb
COPY . .

ARG VERSION=Development
RUN apk add git
RUN go get ./...
RUN go build -ldflags="-X 'github.com/mangaweb4/mangaweb4-backend/system.VersionString=$VERSION' " -o mangaweb4-backend .

# Stage 2 -- build the target image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder1 /go/src/mangaweb/mangaweb4-backend ./

EXPOSE 8972

CMD ["./mangaweb4-backend"]