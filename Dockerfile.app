# Etapa de build
FROM golang:1.20 AS builder

WORKDIR /app

# Copie os arquivos go.mod e go.sum e baixe as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copie o código-fonte
COPY . .

# Compile o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# Etapa final
FROM gcr.io/distroless/base-debian11

WORKDIR /

# Copie o binário compilado da etapa anterior
COPY --from=builder /app/server /server

# Comando de entrada
ENTRYPOINT ["/server"]
