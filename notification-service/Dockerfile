FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /notification-app .

EXPOSE 8081

CMD ["./notification-app"]