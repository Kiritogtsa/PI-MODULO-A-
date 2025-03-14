package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type perguta_resposta struct {
	respota string
	perguta string
}

func (s *perguta_resposta) nome() {
	return
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{}
var exaplechh chan string
var respotas []perguta_resposta

// inicia o websocket
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade error:", err)
		return
	}

	for {
		// ler o chanal, e caso tiver messagem envia para o cliente WebSocket
		select {
		case message := <-exaplechh:
			err = c.WriteMessage(websocket.TextMessage, []byte(message))
			// se acontecer algum erro ele imprime este erro no servidor
			if err != nil {
				log.Println("Write error:", err)
			}

		default:
			log.Println("No message available in channel")
			time.Sleep(time.Second * 2)
		}
	}
}

type Reposta struct {
	Reposta string
}

// ler o respota do arduino e coloca a respota num chanal
func ardcuino(w http.ResponseWriter, r *http.Request) {
	var reposta Reposta
	err := json.NewDecoder(r.Body).Decode(&reposta)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	exaplechh <- reposta.Reposta
}
func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

// função que inicia o servidor
func main() {
	exaplechh = make(chan string)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.HandleFunc("/arduino", ardcuino)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`

<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var ws;

    var print = function(message) {
        var p = document.createElement("p");
        p.textContent = message;
        output.appendChild(p);
        output.scroll(0, output.scrollHeight);
  			console.log(message);
    };

    function connectWebSocket() {
        ws = new WebSocket("{{.}}");
        ws.onmessage = function(evt) {
            print(evt.data); // Exibe apenas a resposta do servidor dentro de um <p>
        };
        ws.onclose = function() {
            setTimeout(connectWebSocket, 3000); // Tenta reconectar após 3s
        };
        ws.onerror = function(evt) {
            console.error("WebSocket Error:", evt);
        };
    }

    connectWebSocket(); // Conecta automaticamente ao carregar a página
});
</script>
</head>
<body>
<div id="output" style="max-height: 70vh; overflow-y: scroll;"></div>
</body>
</html>
`))
