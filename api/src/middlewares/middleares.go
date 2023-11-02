package middlewares

import (
	"api/src/auth"
	"api/src/responses"
	"log"
	"net/http"
)

func Logger (next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requisicao: %s %s", r.Method, r.URL.RequestURI())
		next(w, r)
	}
}

// Chama a função que vai validar o nosso token, verifica se ele é válido e, se for, autoriza o usuário a acessar as funcionalidades do sistema. isso é aplicado antes da rota de login e, se estiver tudo ok, chama a rota de login com o next().
func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := auth.ValidateToken(r); erro != nil {
			responses.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		next(w, r)
	}
}