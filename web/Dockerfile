FROM golang:alpine AS builder

RUN apk update && apk add git && apk add ca-certificates

WORKDIR /go/src/github.com/zale144/instagram-bot/web

COPY . .

RUN set -x && \
	go get github.com/golang/dep/cmd/dep && \
	which dep && \
	date && \
	dep ensure -v && \
	echo "vendor:" && \
	dep status

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o instagram-bot.web .

# second stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/zale144/instagram-bot/web .
