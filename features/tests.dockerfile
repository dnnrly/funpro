FROM golang:1.16

RUN go get -v github.com/cucumber/godog/cmd/godog@v0.11.0

WORKDIR /cache

COPY go.mod .
COPY go.sum .
RUN go mod download

WORKDIR /app

CMD godog -t @Integration