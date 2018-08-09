FROM golang:1.9-alpine3.7 as builder

COPY . /go/src/github.com/flaccid/sesr

WORKDIR /go/src/github.com/flaccid/sesr/cmd/sesr

RUN apk add --update --no-cache git gcc musl-dev && \
    go get ./... && \
    CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /opt/bin/sesr .

FROM centurylink/ca-certs

COPY --from=builder /opt/bin/sesr /opt/bin/sesr

ENTRYPOINT ["/opt/bin/sesr"]
