# Define HEADLESS at the top to be used in FROM instructions
ARG HEADLESS=true

# --- Stage: Real Frontend Build (Active when HEADLESS=false) ---
FROM node:22-alpine AS assets-false
WORKDIR /app/ui
RUN npm install -g pnpm
COPY ui/package.json ui/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY ui/ ./
RUN pnpm run build --outDir ../public

# --- Stage: Empty Assets (Active when HEADLESS=true) ---
FROM alpine:latest AS assets-true
RUN mkdir -p /app/public

# --- Stage: Asset Selector ---
FROM assets-${HEADLESS} AS final-assets

# --- Stage: Build Backend ---
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o toucan cmd/toucan/main.go

# --- Final Production Image ---
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates

# Copy the binary
COPY --from=backend-builder /app/toucan .

# Copy assets from the selected stage
COPY --from=final-assets /app/public ./public

RUN mkdir -p uploads
EXPOSE 8080

# The backend is already configured to serve the public/ folder if it exists
CMD ["./toucan"]
