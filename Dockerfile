FROM golang:alpine as builder
RUN apk add --no-cache git build-base
RUN go get -v github.com/swordbeta/trello-burndown/...

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /root
COPY --from=builder /go/bin/cmd /app/trello-burndown
ENTRYPOINT ["/app/trello-burndown"]

