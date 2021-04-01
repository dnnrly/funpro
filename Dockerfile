FROM golang:1.16-alpine as builder

WORKDIR /cache

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN mkdir /build
WORKDIR /build
COPY . .

RUN go build -o funpro .

FROM golang:1.16-alpine

RUN mkdir /app
COPY --from=builder /build/funpro /app/funpro
WORKDIR /app

CMD ./funpro
