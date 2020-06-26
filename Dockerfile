FROM golang:1.14.4-alpine3.12 AS builder

ADD . /app
WORKDIR /app
RUN go build -o issue-overseer

FROM alpine:3.12.0
WORKDIR /app
COPY --from=builder /app/issue-overseer .
COPY --from=builder  /app/watch.sh .
RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64
RUN chmod +x /usr/local/bin/dumb-init
ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]
CMD ["/app/watch.sh"]
