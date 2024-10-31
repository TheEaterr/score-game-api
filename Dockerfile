FROM alpine:latest

# Install required packages
RUN apk add --no-cache go gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o main .

# Run the application
CMD ["/app/main"]