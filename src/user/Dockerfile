FROM golang:1.17-alpine

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# ENV GIT_TERMINAL_PROMPT=1

# RUN go env -w GOPRIVATE="github.com/startup-of-zero-reais/*"
# RUN go env -w GOPROXY=https://proxy.golang.org/ GOSUMDB=off

ARG github_token

WORKDIR /go/src

RUN apk add --no-cache git
RUN git config --global --add url."https://${github_token}@github.com".insteadOf "https://github.com"

COPY . .

RUN go mod tidy
RUN go build .

ENTRYPOINT [ "echo", "'done'" ]