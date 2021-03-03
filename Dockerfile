FROM golang:1.16.0-alpine3.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

WORKDIR /go/src/app/examples
RUN go build