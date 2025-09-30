# Use a specific Go version and a recent Alpine base for the build stage
FROM golang:1.23-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Build the application as a static binary
RUN CGO_ENABLED=0 go build -o /app/application

# Use a minimal, non-root base image for the final stage
FROM alpine:3.20

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the executable and set ownership
COPY --from=build --chown=appuser:appgroup /app/application /app/executable

# Switch to the non-root user
USER appuser

EXPOSE 8080
ENTRYPOINT ["/app/executable"]
CMD ["server"]
