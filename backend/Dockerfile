# ── Stage 1: BUILD ───────────────────────────────────────────
# Use official Go image to compile the app
FROM golang:1.25-alpine AS builder

# Install git (required for downloading Go modules from GitHub)
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (so Docker caches dependencies)
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of your backend code
COPY . .

# Build the binary — CGO_ENABLED=0 makes it fully static (no external libs needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ── Stage 2: RUN ─────────────────────────────────────────────
# Use a tiny image — final image will be ~20MB instead of ~600MB
FROM alpine:3.19

# Needed for HTTPS calls (e.g. your email sending via gomail)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the compiled binary from the build stage
COPY --from=builder /app/main .

# Expose port 8080 (what your main.go uses)
EXPOSE 8080

# Run the binary
CMD ["./main"]