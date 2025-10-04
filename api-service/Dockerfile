FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api ./cmd/server.go


FROM alpine:latest
COPY  --from=builder /app/api .
EXPOSE 8000
CMD ["./api"]