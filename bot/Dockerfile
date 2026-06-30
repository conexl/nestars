# syntax=docker/dockerfile:1

ARG VITE_DISABLE_AUTH=false

# Stage 1: build Mini App
FROM node:22-alpine AS miniapp-builder
ARG VITE_DISABLE_AUTH
WORKDIR /miniapp
COPY miniapp/package.json miniapp/package-lock.json* ./
RUN npm install
COPY miniapp/ .
RUN VITE_DISABLE_AUTH=${VITE_DISABLE_AUTH} npm run build

# Stage 2: build Go bot
FROM golang:1.26-alpine AS go-builder

WORKDIR /app

COPY bot/go.mod bot/go.sum ./
RUN go mod download

COPY bot/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bot ./cmd/main.go

# Stage 3: runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=go-builder /bot /app/bot
COPY --from=miniapp-builder /miniapp/dist /app/miniapp/dist

ENTRYPOINT ["/app/bot"]
