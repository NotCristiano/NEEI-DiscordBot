# Builder
FROM golang:1.25-alpine AS builder

# Install git
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app

# Copy dependency files first to cache downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# CGO_ENABLED=0 creates a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /neei-bot cmd/neei-discordbot/main.go

# Use a minimal alpine image for the final container
FROM alpine:latest

# Install CA certificates so the bot can talk to Discord (HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the compiled binary from the builder stage
COPY --from=builder /neei-bot .

# Command to run the bot
CMD ["./neei-bot"]