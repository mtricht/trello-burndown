FROM golang:alpine as builder
RUN apk add --no-cache git build-base
RUN go get -v github.com/swordbeta/trello-burndown

FROM alpine
WORKDIR /root
COPY --from=builder /go/bin/trello-burndown /app/trello-burndown
ENTRYPOINT ["/app/trello-burndown"]

