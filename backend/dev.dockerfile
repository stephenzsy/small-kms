FROM golang:1.21

VOLUME /usr/src/app
WORKDIR /usr/src/app

RUN go install github.com/mitranim/gow@latest

EXPOSE 9000

CMD ["gow", "run", "."]