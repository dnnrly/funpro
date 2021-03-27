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

RUN mkdir /build
WORKDIR /build
COPY . .

ENV AWS_DEFAULT_REGION=eu-west-1
ENV AWS_ACCESS_KEY_ID=test
ENV AWS_SECRET_ACCESS_KEY=test

CMD mage -d test install && \
    godog -t @Integration
