FROM golang:1.14.4-alpine3.12 AS builder

ADD . /app
WORKDIR /app
RUN go build -o issue-overseer

FROM alpine:3.12.0
WORKDIR /root
COPY --from=builder /app/issue-overseer .
COPY --from=builder  /app/watch.sh .
CMD sh watch.sh $GITHUB_ORGANIZATION
