FROM golang:1.26.1-alpine AS builder

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./main.go

FROM alpine:3.21

WORKDIR /api

COPY --from=builder /api/server .

RUN adduser -D apiuser
USER apiuser

CMD ["./server"]
