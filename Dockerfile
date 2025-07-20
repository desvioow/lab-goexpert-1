FROM golang:1.22 as build

WORKDIR /app
COPY . .

# Update CA certificates
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Update the CA certificate bundle
RUN update-ca-certificates
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

ENTRYPOINT ["./cloudrun"]