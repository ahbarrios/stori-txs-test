# Build image
FROM golang:1.21-alpine as builder

RUN apk update && apk add openssh git make && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/build/main main.go

# Deploy image
FROM alpine:latest

RUN apk update && apk add ca-certificates tzdata && apk add dnsmasq && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/build/main .

CMD ["./main"]