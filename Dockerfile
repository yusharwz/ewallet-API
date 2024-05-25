FROM golang:alpine

# Update apk and install git and bash
RUN apk update && apk add --no-cache git bash

WORKDIR /app

# Copy the entire project into the container
COPY . .

# Set executable permission for wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

RUN go mod tidy

RUN go build -o final-project-enigma

EXPOSE 8080

# Use wait-for-it.sh to wait for PostgreSQL before starting the application
ENTRYPOINT ["/app/wait-for-it.sh", "postgres:5432", "--", "/app/final-project-enigma"]
