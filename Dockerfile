FROM golang:1.23.3

# Set the working directory inside the container
WORKDIR /cls

# Copy the Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build cmd/cls/server/main.go

# Expose the port your app runs on
EXPOSE 443

RUN mkdir ~/received

# Command to run the application
CMD ["./main"]

