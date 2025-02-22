# First Stage: Build Go binary
FROM golang:1.23.4-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o quizz cmd/server/main.go

# Second Stage: Run binary in a smaller image
FROM alpine:3.21 
WORKDIR /root/
COPY --from=builder /app/quizz .
EXPOSE 8080
CMD ["./quizz"]
