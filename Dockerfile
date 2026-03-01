# Stage 1: Build
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o org-api ./cmd/app/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/org-api .

EXPOSE 8080

CMD ["./org-api"]