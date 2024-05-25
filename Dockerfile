FROM golang:1.16-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o final-project-enigma

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /app/final-project-enigma /app/final-project-enigma
COPY wait-for-it.sh /app/wait-for-it.sh

RUN chmod +x /app/wait-for-it.sh

EXPOSE 8080
ENTRYPOINT ["/app/wait-for-it.sh", "postgres:5432", "--", "/app/final-project-enigma"]
