FROM golang:1.14-alpine

WORKDIR /go/src/twitch-bot
COPY . .

RUN go get -d -v
RUN go install -v -x

CMD ["twitch-bot"]
