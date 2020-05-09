FROM golang:alpine

RUN mkdir /app /temp-build
ADD . /temp-build
WORKDIR /temp-build
RUN go build -o issue-label-bot
RUN mv issue-label-bot ../app
RUN mv watch.sh ../app
WORKDIR /app
RUN rm -rf ../temp-build
CMD sh watch.sh $GITHUB_ORGANIZATION
