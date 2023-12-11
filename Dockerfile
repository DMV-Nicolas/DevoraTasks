# Build stage
FROM golang:1.21.5-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env . 
COPY start.sh .
COPY migrate .
COPY db/migration ./migration

EXPOSE 5000
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]