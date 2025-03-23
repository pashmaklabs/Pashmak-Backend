FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app .

FROM alpine:latest

WORKDIR /app

RUN echo "http://dl-4.alpinelinux.org/alpine/v3.19/main" > /etc/apk/repositories && \
    echo "http://dl-4.alpinelinux.org/alpine/v3.19/community" >> /etc/apk/repositories && \
    apk update && apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/app .

EXPOSE 8080

CMD ["./app"]
