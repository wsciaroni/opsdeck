# Stage 1: Build Frontend
FROM node:22-alpine AS frontend
WORKDIR /app
COPY web/package*.json ./web/
WORKDIR /app/web
RUN npm install
COPY web/ .
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.24-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy frontend build artifacts to expected location for embedding
# We decided to use cmd/server/dist in main.go
COPY --from=frontend /app/web/dist ./cmd/server/dist
RUN go build -o main cmd/server/main.go

# Stage 3: Final Image
FROM alpine:3.19
WORKDIR /app
COPY --from=backend /app/main .
CMD ["./main"]
