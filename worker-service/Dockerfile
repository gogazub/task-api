FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o code_processor ./cmd/main.go


FROM alpine:latest
COPY  --from=builder /app/code_processor .
CMD ["./code_processor"]


# docker run -v /var/run/docker.sock:/var/run/docker.sock code-processor 