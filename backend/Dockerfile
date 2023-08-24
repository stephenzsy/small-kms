FROM golang:1.21 AS build
WORKDIR /go/src
COPY ./* .

ENV CGO_ENABLED=1
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o smallkms .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/smallkms ./
EXPOSE 9000/tcp
ENTRYPOINT ["./smallkms"]
