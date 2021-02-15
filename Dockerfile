FROM golang:alpine

WORKDIR /go/src/github.com/kjcodeacct/layover

COPY . .

RUN apk update && apk add --no-cache git

RUN export CGO_ENABLED=0
RUN export GO111MODULE=on
RUN go get -d -v
RUN go build -o /go/bin/layover

ENTRYPOINT ["/go/bin/layover", "proxy"]