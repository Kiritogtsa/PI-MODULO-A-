from flask import Flask, render_template
import serial
import os
import time

respostas = []
erro = None

def conexao_arduino():
    '''
    Esta função estabelece a comunicação com o arduíno, assim como fica lendo a porta(USB) em que ele esta conectado
    continuadamente pegar os valores dos buttons quando eles são pressionados

    Manti a conexão assim como a leitura junto para não precisar ficar fazendo várias funções
    '''
    global respostas
    porta = os.system('python -m serial.tools.list_ports') # manda um comando para o SO para retornar a porta que o arduíno está conectado
    try:
        arduino = serial.Serial(porta) # Cria o objeto que vai manipular esta porta e passa como argumento a porta que ele vai manipular
        while True: # loop contínuo para ler a porta
            if arduino.in_waiting(): # caso haja algum valor na porta
                respostas.append(arduino.readline()) # armazena esta informação
            time.sleep(0.5)
    except Exception as erro:
        print(erro)

sistema = Flask(__name__,template_folder='./templates')

@sistema.route('/')
def ola():
    return render_template('perguntas.html')

if __name__ == '__main__':
    sistema.run(debug=  True)