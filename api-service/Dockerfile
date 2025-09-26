FROM golang:1.22.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api ./cmd/server.go


FROM alpine:latest
COPY  --from=builder /app/api .
EXPOSE 8000
CMD ["./api"]




# RUN apt-get update && apt-get install -y \
#     python3 \
#     python3-pip \
#     python3-venv \
#     && rm -rf /var/lib/apt/lists/*
#
# RUN ln -s /usr/bin/python3 /usr/bin/python
#
# RUN python3 -m venv /opt/venv
# RUN /opt/venv/bin/pip install pytest requests
#
# ENV PATH="/opt/venv/bin:${PATH}"
#
# WORKDIR /app
#
# COPY . .
#
# RUN go mod download
# RUN go build -o main .
#
# RUN pytest --version
#
# EXPOSE 8000
#
# CMD ./main & sleep 2 && pytest tests.py -v