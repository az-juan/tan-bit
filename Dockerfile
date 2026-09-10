FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# RUN CGO_ENABLED=0 go build -o /app/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /app/api

# Etapa 2: imagen final (solo el binario)
FROM alpine:3.20
RUN adduser -D -u 1000 app
USER app
COPY --from=builder /app/api /usr/local/bin/api
EXPOSE 8080
CMD ["api"]
