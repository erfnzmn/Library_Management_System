FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 

RUN go build -o server ./cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs


EXPOSE 8080
CMD ["./server"]