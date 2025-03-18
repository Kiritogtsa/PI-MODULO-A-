window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

   function print(message) {
        var d = document.getElementById("pergunta");
        d.innerHTML = "<p>" + message + "</p>";     } 
    function connectWebSocket() {
        if (ws) {
            return;
        }
        ws = new WebSocket("ws://localhost:8080/echo");
        ws.onopen = function() {
            console.log("WebSocket conectado.");
        };

        ws.onmessage = function(evt) {
            var data = JSON.parse(evt.data);      
            console.log(data);
            if (data.type == "string") {
                console.log("aqui")
                print(data.msg);

            } else if (data.type === "html") {
                criarCampoInput();
            }
        };

        ws.onclose = function() {
            console.log("WebSocket desconectado.");
            ws = null;
        };

        ws.onerror = function(evt) {
            console.error("Erro no WebSocket:", evt);
        };
    }

    function criarCampoInput() {
        var inputContainer = document.getElementById("input-container");
        inputContainer.innerHTML = '<input type="text" id="nome" placeholder="Digite seu nome">';
        var sendButton = document.createElement("button");
        sendButton.textContent = "Enviar";
        sendButton.onclick = function() {
            var nome = document.getElementById("nome").value;
            if (ws) {
                ws.send(nome);
            }
        };
 
       ws.send(input.value);
        return false;
    };

    

    connectWebSocket(); // Inicia WebSocket ao carregar a página
});

