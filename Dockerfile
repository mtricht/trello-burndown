FROM golang:alpine as builder
RUN apk add --no-cache git build-base
ENV GOBIN=$GOPATH/bin
WORKDIR /src/
COPY . /src/
RUN go install ./...

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /root
COPY --from=builder /go/bin/cmd /app/trello-burndown
ENTRYPOINT ["/app/trello-burndown"]