FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o final-project-enigma

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/final-project-enigma /app/final-project-enigma
COPY .env /app/.env
COPY wait-for-postgres.sh /app/wait-for-postgres.sh

RUN apk add --no-cache netcat-openbsd

EXPOSE 8080

CMD ["/app/wait-for-postgres.sh", "postgres", "5432", "--", "/app/final-project-enigma"]
