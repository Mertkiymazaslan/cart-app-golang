FROM golang:1.21

# Get the environment variables from "docker-compose.yaml"
ARG ENVIRONMENT
ARG DB_URL

# Set the environment variable values
ENV ENVIRONMENT=$ENVIRONMENT
ENV DB_URL=$DB_URL

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY ./ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /checkoutProject ./cmd/

# To bind to a TCP port
EXPOSE 8080

# Run
CMD ["/checkoutProject"]
