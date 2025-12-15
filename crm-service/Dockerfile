# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git untuk go mod yang membutuhkan git
RUN apk add --no-cache git

# Download dependencies dulu (agar caching lebih efisien)
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file source code
COPY . .

# Build binary (static build untuk image yang lebih kecil)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Runtime stage
FROM alpine:latest

# Certificate + timezone
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 3002

CMD ["./main"]
