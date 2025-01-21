FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"] 