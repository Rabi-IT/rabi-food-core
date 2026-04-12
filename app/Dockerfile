FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/server .

RUN adduser -D appuser
USER appuser

CMD ["./server"]
