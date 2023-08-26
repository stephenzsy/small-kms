FROM golang:1.21-alpine as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
COPY common ./common
COPY admin ./admin
COPY migrations ./migrations

RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /smallkms

FROM scratch

# Copy our static executable.
COPY --from=builder /smallkms /smallkms
# Run the hello binary.
EXPOSE 9000/TCP

# Run
CMD ["/smallkms"]
