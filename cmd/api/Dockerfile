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
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /api cmd/api/main.go

# Run stage
FROM alpine:3.19

COPY --from=build /api /app/.env ./

# Expose API port
EXPOSE 3000

CMD ["/api"]
