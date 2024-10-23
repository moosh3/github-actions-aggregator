FROM golang:1.18-alpine

# Install migrate tool
RUN apk add --no-cache curl bash && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | \
    tar xvz && mv migrate.linux-amd64 /usr/local/bin/migrate

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Expose port
EXPOSE 8080

# Run migrations and start the application
CMD ["sh", "-c", "./scripts/migrate.sh up && go run cmd/server/main.go"]
