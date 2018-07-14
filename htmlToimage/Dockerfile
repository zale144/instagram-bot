FROM idocking/wkhtmltopdf:ubuntu

# Install golang and Install go-wkhtmltox
RUN apt-get update \
    && apt-get -y --no-install-recommends install git wget\
    && mkdir -p /tmp/go \
    && cd /tmp/go \
    && wget -q https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz \
    && mkdir -p $HOME/go \
    && export GOPATH=$HOME/go \
    && export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin \
    && cp go get github.com/zale144/instagram-bot/htmlToimage \
    && cd $GOPATH/src/github.com/zale144/instagram-bot/htmlToimage \
    && go build \
    && mkdir -p /app \
    && cp htmlToimage /app \
    && rm -rf /tmp/go/* \
    && rm -rf /usr/local/go \
    && rm -rf $GOPATH \
    && apt-get purge -y --auto-remove git wget

WORKDIR /app
#WORKDIR /go/src/github.com/zale144/instagram-bot/htmlToimage

#COPY . .

#RUN set -x && \
#	go get github.com/golang/dep/cmd/dep && \
#	which dep && \
#	date
#RUN set -x && \
#	dep ensure -v && \
#	echo "vendor:" && \
#	dep status

#UN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o instagram-bot.htmltoimage .

