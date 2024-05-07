# Usamos uma imagem base do Go para compilar nosso código
FROM golang:1.21 AS builder

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia o código fonte para o contêiner
COPY . .

# Compila o código Go em um binário
RUN go build -o main .

# Agora usamos uma imagem mais leve para executar nosso aplicativo compilado
FROM alpine:latest

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia o binário compilado do estágio anterior para o contêiner
COPY --from=builder /app/main .

# Define a porta em que o servidor vai escutar
EXPOSE 8080

# Executa o binário do servidor com os argumentos desejados
CMD ["./main", "httpserver", "-e", "prod"]
