FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o /app/quizz cmd/server/main.go

FROM alpine:3.21

WORKDIR /app
COPY --from=builder /app/quizz /app/quizz

EXPOSE 8080
CMD ["/app/quizz"]