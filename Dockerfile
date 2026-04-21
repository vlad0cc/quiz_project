FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend ./
RUN npm run build

FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod ./
RUN go mod download
COPY . ./
COPY --from=frontend-builder /frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/app ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=backend-builder /tmp/app /app/app
COPY --from=backend-builder /app/migrations /app/migrations
COPY --from=backend-builder /app/frontend/dist /app/frontend/dist
COPY .env /app/.env
EXPOSE 8080
CMD ["/app/app"]
