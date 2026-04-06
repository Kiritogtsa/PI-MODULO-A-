package dto

import (
	"database/sql"
	"fmt"

	"github.com/kiritogtsa/PI-MODULO-A-/websocket/model"
)

type Repostasdb struct {
	DB *sql.DB
}

func (r *Repostasdb) Getbyid(id int) (*model.PerguntaResposta, error) {
	query := "select * from pergunta where id_pergunta=?"
	row := r.DB.QueryRow(query, id)
	var pergunta model.PerguntaResposta
	err := row.Scan(&pergunta.ID, &pergunta.Pergunta)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("nenhuma pergunta encontrada para o ID %d", id)
		}
		return nil, err
	}
	return &pergunta, nil
}
