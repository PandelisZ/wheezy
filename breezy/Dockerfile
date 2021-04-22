FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /build
ENV GO111MODULE=on \
    CGO_ENABLED=1
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./breezy .
EXPOSE 8080
CMD ["/build/breezy"]
