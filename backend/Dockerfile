FROM golang:1.21-alpine as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.* ./

RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o server
RUN CGO_ENABLED=0 GOOS=linux go build -o agent-server ./agent

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

# Copy our static executable.
COPY --from=builder /app/server /server
COPY --from=builder /app/agent-server /agnet-server

# Run
CMD ["/server", "admin", ":9000"]
