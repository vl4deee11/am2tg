# Context dir .
FROM golang:1.15-alpine as build

WORKDIR /go/src/am2tg
COPY . /go/src/am2tg

RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -o am2tg

CMD ["./am2tg"]
