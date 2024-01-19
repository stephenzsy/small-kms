FROM golang:1.21.6-alpine as builder

RUN apk add --no-cache gcc g++

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.* ./

RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./
ARG BUILD_ID=dev

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o server -ldflags="-X 'main.BuildID=${BUILD_ID}'" .
RUN CGO_ENABLED=0 GOOS=linux go build -o agent-server -ldflags="-X 'main.BuildID=${BUILD_ID}'" ./agent

FROM alpine:3.19.0

RUN apk add --no-cache ca-certificates openssl

# Copy our static executable.
COPY --from=builder /app/server /server
COPY --from=builder /app/agent-server /agent-server

# Run
CMD ["/server", "admin", ":9000"]
