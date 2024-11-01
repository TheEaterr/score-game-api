FROM alpine:latest

# Install required packages
RUN apk add --no-cache go gcc musl-dev
RUN go mod download github.com/mattn/go-sqlite3@latest

# Set the working directory
WORKDIR /app

EXPOSE 8080

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o main .

# Run the application
CMD ["/app/main"]