# --- Builder Stage ---
FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN mkdir -p /go/pkg/mod /root/.cache/go-build

# Copy go.mod and go.sum first to leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
# Note: We copy the whole backend directory context here
COPY . .

# Build the server application specifically
WORKDIR /build/services/server/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /server_app .

# Build other services if needed in the future
# WORKDIR /build/services/faces/cmd
# RUN CGO_ENABLED=0 GOOS=linux go build -v -o /faces_app .


# --- Final Stage ---
FROM alpine:latest

WORKDIR /app

# Copy only the built application binary from the builder stage
COPY --from=builder /build/db/schema.sql .
COPY --from=builder /server_app .

# Copy other built services if needed
# COPY --from=builder /faces_app .

# Copy the env file for reference/defaults (will be overridden by compose)
# COPY services/server/cmd/app.env .

# Expose the port the Go application listens on (from your config)
EXPOSE 8080

# Command to run the application
# We expect configuration via environment variables passed by Docker Compose
CMD [ "./server_app" ]
