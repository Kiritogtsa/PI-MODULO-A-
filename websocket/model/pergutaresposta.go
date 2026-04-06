package model

type PerguntaResposta struct {
	ID       int    `json:"id"`
	Resposta string `json:"resposta"`
	Pergunta string `json:"pergunta"`
}
