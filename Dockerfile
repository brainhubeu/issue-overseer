FROM golang:1.14.4-alpine3.12

RUN mkdir /app /temp-build
ADD . /temp-build
WORKDIR /temp-build
RUN go build -o issue-overseer
RUN mv issue-overseer ../app
RUN mv watch.sh ../app
WORKDIR /app
RUN rm -rf ../temp-build
CMD sh watch.sh $GITHUB_ORGANIZATION
