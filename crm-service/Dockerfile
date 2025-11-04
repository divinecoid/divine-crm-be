# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o divine-crm cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser

WORKDIR /root/

# Copy binary
COPY --from=builder /app/divine-crm .

# Copy migrations if exists (FIXED!)
# Option 1: Skip if not exists
RUN mkdir -p ./migrations

# Option 2: Copy only if exists (comment out line 39-40 if using this)
# COPY --from=builder /app/migrations ./migrations 2>/dev/null || true

# Change ownership
RUN chown -R appuser:appuser /root/

USER appuser

EXPOSE 3002

CMD ["./divine-crm"]