FROM golang:1.21

COPY go.mod go.sum /go/src/
WORKDIR /go/src

COPY . /go/src/

RUN go mod download

CMD go test github.com/r4chi7/aspire-lite/e2e
