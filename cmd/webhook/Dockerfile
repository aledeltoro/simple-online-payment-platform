# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build

# Create work directory
WORKDIR /app

# Copy and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy API packages
COPY . ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /webhooks cmd/webhook/main.go

# Run stage
FROM alpine:3.19

COPY --from=build /webhooks .

# Expose API port
EXPOSE 3001

CMD ["/webhooks"]