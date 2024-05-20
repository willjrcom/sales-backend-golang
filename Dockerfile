# Use uma imagem base do Go
FROM golang:1.21

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código-fonte para dentro do contêiner
COPY . .

RUN brew install librdkafka

# Execute o servidor diretamente
CMD ["go", "run", "main.go", "httpserver", "-e", "prod"]
