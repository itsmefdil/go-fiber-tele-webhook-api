FROM golang:1.23.1

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 3000

# Set environment variables (optional)
# ENV BASIC_AUTH_USERNAME=admin
# ENV BASIC_AUTH_PASSWORD=supersecret

# Run the application
CMD ["./main"]
