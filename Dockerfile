FROM golang:alpine

ADD *.go /app/
ADD watch.sh /app/
WORKDIR /app
CMD sh watch.sh $GITHUB_ORGANIZATION
