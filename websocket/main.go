package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type repostasdb struct {
	DB *sql.DB
}

func (r *repostasdb) getbyid(id int) (*PerguntaResposta, error) {
	query := "select * from pergunta where id_pergunta=?"
	row := r.DB.QueryRow(query, id)
	var pergunta PerguntaResposta
	err := row.Scan(&pergunta.ID, &pergunta.Pergunta)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("nenhuma pergunta encontrada para o ID %d", id)
		}
		return nil, err
	}
	return &pergunta, nil
}

type PerguntaResposta struct {
	ID       int    `json:"id"`
	Resposta string `json:"resposta"`
	Pergunta string `json:"pergunta"`
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{}
var exaplechh chan string
var respostas []PerguntaResposta

func (s* PerguntaResposta) associar_id_pergunta_resposta(id int,pergunta string, resposta string){
	// Crio o objeto resposta
    pr := PerguntaResposta{
        ID:       id,
        Pergunta: pergunta,
        Resposta: resposta,
    }

    // Adiciona a nova resposta à lista de respostas
    respostas = append(respostas, pr)

    // Limita o número de respostas a 8 que é a quantidade de perguntas no banco de dados
    if len(respostas) > 8 {
        return
    }

	// chama a função para buscar a pergunta pelo próximo id de acordo com o próximo id d lista
	db: = repostasdb{}
	var nova_pergunta = db.getbyid(len(PerguntaResposta)+1)
	return nova_pergunta
}

// inicia o websocket
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	var id int = 1
	if err != nil {
		log.Print("Upgrade error:", err)
	}
	repdb, err := abrir_banco_dados()
	if err != nil {
		log.Print("Upgrade error:", err)
		return
	}
	var perg *PerguntaResposta
	if id > len(respostas)-1 {
		perg, err = repdb.getbyid(id)
		log.Println(perg, err)
		if err != nil {
			log.Println(err)
			id = 1
			log.Println(id)
			log.Println(respostas)
		} else {
			respostas = append(respostas, *perg)
		}
	} else {
		perg = &respostas[id-1]
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(perg.Pergunta))
	if err != nil {
		log.Println("Write error:", err)
	}
	id++
	for {
		// ler o chanal, e caso tiver messagem envia para o cliente WebSocket
		select {
		case message := <-exaplechh:
			var perg *PerguntaResposta
			log.Println(id)
			log.Println(len(respostas))
			if id < len(respostas)-1 {
				perg = &respostas[id-1]
			} else {
				perg, err = repdb.getbyid(id)
				log.Println(err)
				if err != nil {
					log.Println("teste", err)
				} else {
					respostas = append(respostas, *perg)
				}
			}
			err = c.WriteMessage(websocket.TextMessage, []byte(perg.Pergunta))
			if err != nil {
				log.Println("Write error:", err)
			}
			respostas[id-1].Resposta = message
			id++
			if id == 9 {
				id = 1
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
func abrir_banco_dados() (*repostasdb, error) {

	db, err := sql.Open("sqlite3", "./pibancodedados")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("aqui")
	return &repostasdb{DB: db}, nil
}

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
