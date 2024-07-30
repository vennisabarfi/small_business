# Build image with golang:latest
FROM golang:latest as build

# Set working directory inside container
WORKDIR /small_business

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Download go module dependencies
RUN go mod download

# Copy source code into container
COPY . .

# Compile the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go_backend 

# Final stage with alpine:latest
FROM alpine:latest as run

# Set working directory inside container
WORKDIR /small_business

EXPOSE 3000

# Copy the compiled binary from the builder stage
COPY --from=build /go_backend /go_backend

# Command to run the executable
CMD ["/go_backend"]
