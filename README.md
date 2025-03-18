# PI-MODULO-A-

windows

instalar o python3 e pip
as dependecias do python3
pyserial, flask e async

linux

criar um ambiente python3 virtual com o comando python -m venv nome,
depois ativar o ambiente virtual com o nome/bin/active,
depois para sair e so deactivave
no ambiente virtual instalar com o pip as seguites blibioteclas,
fask, pyserial e async
Deve criar um arquivo JSON na pasta `websocket` para que o sistema possa enviar e-mails.

**Exemplo de arquivo `variaveis.json`:**
```json
{
  "email": "seu-email@gmail.com",
  "email_send": "destinatario@gmail.com",
  "key": "sua-chave-de-autenticacao",
  "api": "smtp.seu-servidor.com",
  "port": 587
}
