# Context dir .
FROM golang:1.15 as build

WORKDIR /go/src/am2tg
COPY . /go/src/am2tg

RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -o am2tg

ENTRYPOINT ["./am2tg"]
