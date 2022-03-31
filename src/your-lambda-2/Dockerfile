FROM golang:1.17-alpine

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /go/src

COPY . .

RUN go mod tidy
RUN go build .

ENTRYPOINT [ "echo", "'done'" ]