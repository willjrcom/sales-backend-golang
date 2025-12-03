# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copie os arquivos de dependência primeiro (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copie o código-fonte
COPY . .

# Compile o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sales-backend main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Instale certificados SSL (necessário para HTTPS)
RUN apk --no-cache add ca-certificates tzdata

# Copie o binário compilado
COPY --from=builder /app/sales-backend .

# Copie os arquivos de migração (caso precise rodar migrate)
COPY --from=builder /app/bootstrap/database/migrations ./bootstrap/database/migrations

# Copie o script de entrypoint
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Expose a porta padrão
EXPOSE 8080

# Use o entrypoint script para migrações automáticas
# O container executa migrate-all automaticamente antes do servidor
# Para desabilitar: defina AUTO_MIGRATE=false
ENTRYPOINT ["/app/docker-entrypoint.sh"]

# Comando padrão: httpserver em produção
CMD ["httpserver", "-e", "prod"]
