FROM golang:alpine AS builder

RUN apk update && apk add git && apk add ca-certificates

WORKDIR /go/src/instagram-bot/web
COPY . .
RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o instagram-bot.web .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/instagram-bot.web .
