# ---- Build stage ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache nodejs npm

WORKDIR /app

# Cache Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Cache npm dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy source
COPY . .

# Generate templ files
RUN go tool templ generate

# Build Tailwind CSS
RUN npx @tailwindcss/cli \
    -i ./static/css/input.css \
    -o ./static/css/output.css \
    --minify

# Install goose for running migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build

# ---- Runtime stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy binary and static assets from builder
COPY --from=builder /app/fithub .
COPY --from=builder /app/static ./static

# Copy goose binary and migrations for release_command
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/sql/schema ./sql/schema

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/v1/healthz || exit 1

CMD ["./fithub"]
