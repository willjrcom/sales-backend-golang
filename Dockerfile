# Use uma imagem base do Go
FROM golang:1.23

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /go/app

# Copie o código-fonte para dentro do contêiner
COPY . .

# RUN apt-get update && apt-get install -y librdkafka-dev

# Execute o servidor diretamente
# "-e", "prod"
CMD go run main.go httpserver

# Video full cycle kafka
# CMD ["tail", "-f", "/dev/null"]