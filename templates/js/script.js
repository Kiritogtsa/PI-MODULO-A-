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
    function removercampopertunta() {
        var pergunta = document.getElementById("pergunta");
        pergunta.innerHTML = "";
    }
    function removeCampoInput() {
        var inputContainer = document.getElementById("input-container");
        inputContainer.innerHTML = "";
    }
    function criarCampoInput() { 
        var inputContainer = document.getElementById("input-container");
        inputContainer.innerHTML = '<input type="text" id="nome" placeholder="Digite seu nome">';
        var sendButton = document.createElement("button");
        sendButton.textContent = "Enviar";
        sendButton.onclick = function() {
            var nome = document.getElementById("nome").value;
            if (ws) {
                console.log("Enviando nome:", nome);
                ws.send(nome);
                removeCampoInput();
            }
        };
        inputContainer.appendChild(sendButton);
        document.getElementById("nome").addEventListener("keyup", function(event) {
            if (event.keyCode === 13) {
                ws.send(document.getElementById("nome").value);
            }
        }); 
    };

    

    connectWebSocket(); // Inicia WebSocket ao carregar a página
});

