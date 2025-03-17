function hover(botao){
    if(botao == 'sim'){
        let botao_sim = document.getElementById('sim');
        botao_sim.style.border = '1px solid black';
    }else if(botao == 'nao'){
        let botao_nao = document.getElementById('nao');
        botao_sim.style.border = '1px solid black';
    }else if(botao == 'nao_sei'){
        let botao_nao_sei = document.getElementById('nao_sei');
        botao_sim.style.border = '1px solid black';
    }
}

function desativa_hover(botao){
    if(botao == 'sim'){
        let botao_sim = document.getElementById('sim');
        botao_sim.style.border = 'none';
    }else if(botao == 'nao'){
        let botao_nao = document.getElementById('nao');
        botao_sim.style.border = 'none';
    }else if(botao == 'nao_sei'){
        let botao_nao_sei = document.getElementById('nao_sei');
        botao_sim.style.border = 'none';
    }
}

window.addEventListener("load", function(evt) {
var output = document.getElementById("output");
var input = document.getElementById("input");
var ws;

var print = function(message) {
    var d = document.createElement("div");
    d.textContent = message;
    output.appendChild(d);
    output.scroll(0, output.scrollHeight);
};
// uma função para remover a os conteudos atuals da tela
// uma função para colocar um campo de input na tela
// e depois enviar via ws o nome que vai vim do input
// uma função para recriar o contudo origal da pagina

// Abre a conexão WebSocket automaticamente
function connectWebSocket() {
    if (ws) {
        return;
    }
    ws = new WebSocket("ws://localhost:8080/echo");
    ws.onmessage = function(evt) {
        if (evt.data){
        print(evt.data); // Exibe apenas a resposta do servidor
        }
    };
    ws.onclose = function() {
        ws = null;
    };
    ws.onerror = function(evt) {
        console.error("WebSocket Error:", evt);
    };
}

connectWebSocket(); // Chama a função automaticamente ao carregar a página

document.getElementById("send").onclick = function(evt) {
    if (!ws) {
        return false;
    }
    ws.send(input.value);
    return false;
};

document.getElementById("close").onclick = function(evt) {
    if (ws) {
        ws.close();
    }
    return false;
};

});