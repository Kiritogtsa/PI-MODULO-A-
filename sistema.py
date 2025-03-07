from flask import Flask, render_template
import serial
import os
import time

respostas = []
erro = None
porta = "/dev/ttyUSB0"
baud_rate = 9000


def conexao_arduino():
    """
    Esta função estabelece a comunicação com o arduíno, assim como fica lendo a porta(USB) em que ele esta conectado
    continuadamente pegar os valores dos buttons quando eles são pressionados

    Manti a conexão assim como a leitura junto para não precisar ficar fazendo várias funções
    """
    global respostas
    # porta = os.system('python -m serial.tools.list_ports') # manda um comando para o SO para retornar a porta que o arduíno está conectado, nao retorna a
    # porta corretam retorna um object int
    # a conexao com o arduino nao pode tar dentro do try_catch, sei la o porque, so não pode pq decodifica errado, nao do padrao utf-8
    arduino = serial.Serial("/dev/ttyUSB0", 9600, timeout=1)
    print("Conectado ao Arduino na porta /dev/ttyUSB0")

    while True:
        if arduino.in_waiting > 0:  # Se houver dados na porta serial
            dados = arduino.readline()  # Lê uma linha completa da serial
            try:
                # Tenta decodificar com UTF-8
                dados_decodificados = dados.decode("utf-8").strip()  # UTF-8
                print(f"Recebido: {dados_decodificados}")
                respostas = respostas.append(dados_decodificados)
            finally:
                print("")


sistema = Flask(__name__, template_folder="./templates")


@sistema.route("/")
def ola():
    return render_template("perguntas.html")


if __name__ == "__main__":
    sistema.run(debug=True)
conexao_arduino()
