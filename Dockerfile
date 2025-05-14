# Build stage
FROM golang:1.24 AS builder
RUN apt-get update && apt-get install -y \
    libwebp-dev \
    gcc \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -v -o /usr/local/bin/app .

# Final stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libwebp7 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /usr/local/bin/app .
EXPOSE 8080
CMD ["./app"]