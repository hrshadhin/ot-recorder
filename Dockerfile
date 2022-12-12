ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine AS builder

# Install dependencies
RUN apk add --no-cache git gcc musl-dev && rm -rf /var/cache/apk/*

# Working directory outside $GOPATH
WORKDIR /build

# Copy go module files and download dependencies
COPY go.* ./
RUN go mod download

# Add source files
ADD . .

# Build the Go app
ARG BUILD_VERSION=0.0.1
ARG BUILD_TIME=00000000-000000
RUN go generate ./cmd && GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X ot-recorder/app.Version=$BUILD_VERSION -X ot-recorder/app.BuildTime=$BUILD_TIME" -o ot-recorder .

# Minimal image for running the application
FROM alpine as final

LABEL org.opencontainers.image.source="https://github.com/hrshadhin/ot-recoder" \
      org.opencontainers.image.url="https://github.com/hrshadhin/ot-recoder" \
      org.opencontainers.image.title="A note taking applications"

# Install/Create dependent tools,location,directory
RUN apk add --no-cache curl tini && \
    rm -rf /var/cache/apk/* && \
    mkdir /persist && chown -R 1000:1000 /persist

# Import the compiled executable from the first stage.
COPY --from=builder /build/ot-recorder /app/ot-recorder

## Open port
EXPOSE 8000

## Perform any further action as an unprivileged location.
USER 1000
WORKDIR /app

HEALTHCHECK --interval=1m --timeout=1s --retries=3 --start-period=2s CMD ["curl", "--fail", "http://localhost:8000/health"]

## Run the compiled binary.
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/ot-recorder","--config","/app/config.yml", "serve"]
