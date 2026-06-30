# syntax=docker/dockerfile:1

ARG VITE_DISABLE_AUTH=false

# Stage 1: build Mini App static files
FROM node:22-alpine AS miniapp-builder
ARG VITE_DISABLE_AUTH
WORKDIR /miniapp
COPY miniapp/package.json miniapp/package-lock.json* ./
RUN npm install
COPY miniapp/ .
RUN VITE_DISABLE_AUTH=${VITE_DISABLE_AUTH} npm run build

# Stage 2: build Go binaries
FROM golang:1.26-alpine AS go-builder
WORKDIR /app
COPY bot/go.mod bot/go.sum ./
RUN go mod download
COPY bot/ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bot ./cmd/bot/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api/main.go
# Legacy full build (bot + API + miniapp in one process)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /legacy ./cmd/main.go

# Stage 3a: Telegram bot runtime (no HTTP API, no miniapp)
FROM alpine:latest AS bot
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=go-builder /bot /app/bot
ENTRYPOINT ["/app/bot"]

# Stage 3b: API + miniapp runtime (no Telegram bot)
FROM alpine:latest AS api
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=go-builder /api /app/api
COPY --from=miniapp-builder /miniapp/dist /app/miniapp/dist
ENTRYPOINT ["/app/api"]

# Stage 3c: legacy full runtime
FROM alpine:latest AS legacy
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=go-builder /legacy /app/legacy
COPY --from=miniapp-builder /miniapp/dist /app/miniapp/dist
ENTRYPOINT ["/app/legacy"]
