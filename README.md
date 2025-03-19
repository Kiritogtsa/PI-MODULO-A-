# PI-MODULO-A-

Este projeto é uma aplicação que integra um Arduino com um servidor WebSocket para coletar e processar respostas de usuários.

## Estrutura do Projeto

- `arduino/`: Código Go para comunicação com o Arduino.
- `cliente_teste/`: Cliente de teste em Go para simular respostas do Arduino.
- `scriptarduino/`: Código Arduino para leitura de botões.
- `websocket/`: Servidor WebSocket em Go para processar e armazenar respostas.

## Requisitos

### Docker

1. Instale o Docker seguindo as instruções oficiais: [Instalação do Docker](https://docs.docker.com/get-docker/).
2. Verifique se o Docker está instalado corretamente:
   ```sh
   docker --version
   ```

### Configuração do Servidor de E-mail

Crie um arquivo JSON na pasta `websocket` para que o sistema possa enviar e-mails.

**Exemplo de arquivo `variaveis.json`:**
```json
{
  "email": "seu-email@gmail.com",
  "email_send": "destinatario@gmail.com",
  "key": "sua-chave-de-autenticacao",
  "api": "smtp.seu-servidor.com",
  "port": 587
}
```

## Executando o Projeto

### Arduino

1. Carregue o código do Arduino localizado em `scriptarduino/scriptdoarduino/scriptdoarduino.ino` no seu Arduino.

### Servidor WebSocket

1. Navegue até a pasta `websocket`:
   ```sh
   cd websocket
   ```
2. Construa a imagem Docker:
   ```sh
   docker build -t minha-imagem .
   ```
3. Execute o container Docker:
   ```sh
   docker run -d -p 8080:8080 --name meu-servidor minha-imagem
   ```

### Cliente de Teste

1. Navegue até a pasta `cliente_teste`:
   ```sh
   cd cliente_teste
   ```
2. Execute o cliente de teste:
   ```sh
   go run main.go
   ```

## .gitignore

O arquivo `.gitignore` está configurado para ignorar os seguintes itens:
- `python_teste`
- `__pycache__`
- `.mypy_cache`
- `websocket/pibancodedados`
- `websocket/log/`
- `websocket/*.json`
