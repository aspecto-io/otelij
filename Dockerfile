FROM golang:1.19-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY config ./config
COPY exporter ./exporter
COPY generator ./generator

COPY *.go ./

RUN go build -o /main

CMD ["/main"]
