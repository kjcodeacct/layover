FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/kjcodeacct/layover/src
COPY . .

RUN export CGO_ENABLED=0 && go build -o /go/bin/layover

ENTRYPOINT ["/go/bin/layover"]