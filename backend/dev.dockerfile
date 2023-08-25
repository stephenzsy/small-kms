FROM golang:1.21

VOLUME /usr/src/app
WORKDIR /usr/src/app

EXPOSE 9000

CMD ["gow", "run", "."]