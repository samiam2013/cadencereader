# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./apps/web/main

# Run stage
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/app .
COPY apps/web/main/views /views
COPY database/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["./app"]

