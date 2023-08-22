FROM golang:1.21

VOLUME /usr/src/app
WORKDIR /usr/src/app

CMD ["go", "run", "./main.go"]