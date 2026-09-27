FROM golang:1.27-alpine AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/app /app/app
EXPOSE 8080
CMD ["/app/app"]
