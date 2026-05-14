FROM golang:1.26.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app_bin ./cmd/api/

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app_bin ./server
COPY --from=builder /build/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
