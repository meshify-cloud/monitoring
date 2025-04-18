# Build stage
FROM golang:1.21-alpine3.19 AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY . .

# Descargar dependencias
RUN go mod download

# Compilar la aplicación
RUN CGO_ENABLED=1 GOOS=linux go build -o main main.go

# Etapa final
FROM alpine:3.19

# Instalar SQLite y otras dependencias necesarias
RUN apk add --no-cache sqlite-libs docker-cli

WORKDIR /app

# Copiar el binario compilado y el archivo monitor.go
COPY --from=builder /app/main ./main
COPY --from=builder /app/main.go ./monitor.go

# COPY --from=builder /app/apps/golang/.env ./.env

# Exponer el puerto
ENV PORT=3001
EXPOSE 3001

# Ejecutar la aplicación
CMD ["./main"]