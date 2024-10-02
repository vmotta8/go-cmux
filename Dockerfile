# Use uma imagem base do Go
FROM golang:1.20

WORKDIR /app

# Copie todo o código-fonte
COPY . .

# Compile o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# Comando de entrada
ENTRYPOINT ["./server"]
