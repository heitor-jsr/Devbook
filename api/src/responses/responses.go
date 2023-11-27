package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

// para retornar uma resposta em JSON, usamos a função abaixo. o parametro de dados precisa ser uma interface generica para ser reaproveitada por outras funções.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		// esse encode que vai transformar os dados que são passados em json.
		if erro := json.NewEncoder(w).Encode(data); erro != nil {
			log.Fatal(erro)
		}
	}

}

func Erro(w http.ResponseWriter, statusCode int, erro error) {
	JSON(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})

}
