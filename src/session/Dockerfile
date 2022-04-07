FROM golang:1.17-alpine

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

ARG github_token

WORKDIR /go/src

RUN apk add --no-cache git
RUN git config --global --add url."https://${github_token}@github.com".insteadOf "https://github.com"

COPY . .

RUN go mod tidy
RUN go build .
