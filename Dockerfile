# Use a specific Go version and a recent Alpine base for the build stage
FROM golang:1.22-alpine AS build

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
FROM alpine:latest
COPY --from=build /app/application /app/executable
EXPOSE 8080
ENTRYPOINT ["/app/executable"]
