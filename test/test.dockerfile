FROM golang:1.16

WORKDIR /cache

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN apt-get update \
    && apt-get install -y zip git python3 python3-pip \
    && pip3 install --upgrade pip \
    && pip3 install awscli \
    && git clone https://github.com/magefile/mage \
    && cd mage \
    && go run bootstrap.go \
    && go get github.com/cucumber/godog/cmd/godog@v0.11.0

RUN mkdir /test
WORKDIR /test
COPY go.mod .
COPY go.sum .
COPY test/ . 

CMD mage install \
    && cd /app \
    && godog -t @Integration
