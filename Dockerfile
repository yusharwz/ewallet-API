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

EXPOSE 8080

CMD ["/app/final-project-enigma"]
