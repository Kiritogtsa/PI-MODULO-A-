package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/mail.v2"
)

type repostasdb struct {
	DB *sql.DB
}

type PerguntaResposta struct {
	ID       int    `json:"id"`
	Resposta string `json:"resposta"`
	Pergunta string `json:"pergunta"`
}

type Data struct {
	Tipo string `json:"type"`
	Msg  string `json:"msg"`
	Res  string `json:"resp"`
}

type Data_log struct {
	Name string
	Date *[]PerguntaResposta
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

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var upgrader = websocket.Upgrader{}
var exaplechh chan string
var respostas = make([]PerguntaResposta, 0)

func perguta_inicial(db *repostasdb) (*PerguntaResposta, int) {
	perg, err := db.getbyid(1)
	if err != nil {
		log.Printf("deu algum erro ao receber a pergunta: %v", err)
		return nil, 0
	}
	return perg, 2
}

func associar_id_pergunta_resposta(id int, pergunta string, resposta string, db *repostasdb) (*PerguntaResposta, int, bool) {
	// Crio o objeto resposta
	pr := PerguntaResposta{
		ID:       id,
		Pergunta: pergunta,
		Resposta: resposta,
	}
	// a gente ainda pode usar essa funçao para verificar quantas vezes o usuario respodeu nao_sei
	// e com isso a gente pode setar o terminou como true e parar este usuario

	// variavel para verificar se as perguntas terminaram
	var terminou bool
	// verifica no banco de dados se tem mais uma perguta
	nova_pergunta, err := db.getbyid(id)
	// se dar um erro entao nao existe mais pergunta, entao ele so fica reatribuindo valores para o
	fmt.Println(len(respostas))
	if err != nil {
		// reseta o id renferenciado, para 1
		id = 1
		terminou = true
		nova_pergunta, _ = db.getbyid(id)
	} else if len(respostas) < id {
		// aqui ele verifica se o tamanho do array e menor que o id, se for ele vai adicionar o novo objeto no array
		respostas = append(respostas, pr)
		id++
	} else {
		// se o tamanho do array for maior que o id - 1, ele nao adiciona mais valores, e sim recoloca eles
		respostas[id-1] = pr
		id++
	}
	// retorna a nova pergunta, sem resposta
	return nova_pergunta, id, terminou
}

func enviar_message(msg []byte, ws *websocket.Conn) error {
	fmt.Println(string(msg))
	return ws.WriteMessage(websocket.TextMessage, msg)
}

func covertdatatojson(tipo string, p *PerguntaResposta, s string) ([]byte, error) {
	var d Data
	if p == nil {
		d = Data{
			Tipo: tipo,
			Msg:  "Nenhuma pergunta disponível", // Evita acessar um ponteiro nulo
			Res:  s,
		}
	} else {
		d = Data{
			Tipo: tipo,
			Msg:  p.Pergunta,
			Res:  s,
		}
	}
	jsondata, err := json.Marshal(d)
	if err != nil {
		return []byte(""), nil
	}
	return jsondata, nil
}

// função para verificar se a pasta dos logs ja existe, se nao existir ele cria
func isexist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if err := os.Mkdir(path, 0777); err != nil {
			return false
		}
		return true
	}
	return true
}
func save_log(path, nome string) error {
	return nil
}

// funçao para criar o arquvo com as repostas
// tambem precisa verificar se e para colocar no banco de dados como um log, dai eu preciso pegar o path
func criar_o_arquivo(d Data_log) error {
	log.Println(isexist("log"))
	data, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}
	path := "log/" + d.Name + time.DateTime + ".txt"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}
	file.Write(data)
	sendemail(path, d.Name)
	return nil
}
func sendemail(path, nome string) {
	arquivo, err := os.Open("./variaveis.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	fmt.Println(arquivo)
	defer arquivo.Close()
	type jsondata struct {
		Email      string `json:"email"`
		Email_send string `json:"email_send"`
		Key        string `json:"key"`
		Api        string `json:"api"`
		Port       int    `json:"port"`
	}
	var data jsondata
	err = json.NewDecoder(arquivo).Decode(&data)
	if err != nil {
		log.Println("sei la: ", err)
	}
	fmt.Println(data)
	message := mail.NewMessage()
	message.SetHeader("From", data.Email)
	message.SetHeader("To", data.Email_send)
	message.SetHeader("Subject", "Log do cliente: "+nome)
	message.SetBody("text/html",
		"<html>"+
			"<body>"+
			"<h1>Log do cliente com o nome :"+nome+"</h1>"+
			"</body>"+
			"</html>")
	// aqui a gente vai anexar o arquivo
	message.Attach(path)
	dialer := mail.NewDialer(data.Api, data.Port, data.Email, data.Key)
	if err := dialer.DialAndSend(message); err != nil {
		log.Println(err)
	}
}

// inicia o websocket
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade error:", err)
	}
	fmt.Println("aqui")
	repdb, err := abrir_banco_dados()
	if err != nil {
		log.Print("Upgrade error:", err)
		return
	}
	perg, id := perguta_inicial(repdb)
	jsondata, err := covertdatatojson("string", perg, "")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(jsondata)
	log.Println(enviar_message(jsondata, c))
	id++
	var terminou bool = false
	for {
		// ler o chanal, e caso tiver messagem envia para o cliente WebSocket
		select {
		case message := <-exaplechh:

			perg, id, terminou = associar_id_pergunta_resposta(id, perg.Pergunta, message, repdb)
			jsondata, err := covertdatatojson("string", perg, message)
			if err != nil {
				log.Println(err)
			}
			log.Println(enviar_message(jsondata, c))
			if terminou {
				fmt.Println("aqui")
				terminou = false
				jsondata, err := covertdatatojson("html", nil, "sei la depois a gente pensa melhor nisso")
				if err != nil {
					log.Println(err)
				}
				log.Println(enviar_message(jsondata, c))
				_, message, err := c.ReadMessage()
				log.Println(string(message))
				if err != nil {
					log.Println(err)
				}
				file_conteudo := Data_log{
					Name: string(message),
					Date: &respostas,
				}
				log.Println(criar_o_arquivo(file_conteudo))
			}
		default:
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
	http.ServeFile(w, r, "templates/index.html")
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
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.HandleFunc("/arduino", ardcuino)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
