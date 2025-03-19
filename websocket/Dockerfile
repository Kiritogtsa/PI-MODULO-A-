
# Imagem base do Golang com Alpine (mais leve e eficiente)
FROM golang:1.24.1-alpine

# Definir diretório de trabalho
WORKDIR /app

# Instalar dependências do sistema necessárias para SQLite e build
RUN apk add --no-cache sqlite-libs gcc musl-dev

# Copiar os arquivos do projeto
COPY . .

# Garantir que CGO esteja ativado para compilar com SQLite
ENV CGO_ENABLED=1

# Baixar dependências do Go
RUN go mod tidy

# Compilar a aplicação
RUN go build -o server .

# Expor a porta usada pelo servidor
EXPOSE 8080

# Comando para executar o servidor
CMD ["./server"]

