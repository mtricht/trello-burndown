FROM alpine:3.4

RUN apk --no-cache add ca-certificates

# Credits to MailHog
RUN apk --no-cache add --virtual build-dependencies go git build-base \
  && mkdir -p /root/gocode \
  && export GOPATH=/root/gocode \
  && go get github.com/swordbeta/trello-burndown \
  && mv /root/gocode/bin/trello-burndown /usr/local/bin \
  && rm -rf /root/gocode \
  && apk del --purge build-dependencies build-base

WORKDIR /root
ENTRYPOINT ["trello-burndown"]
