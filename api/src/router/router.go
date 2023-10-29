package router

import "github.com/gorilla/mux"

// vai gerar e retornar um ponteiro para o router. com isso, retornamos um novo router, com as rotas já configuradas para serem utilizadas nos outros pacotes.
func Generator() *mux.Router{
	return mux.NewRouter()
}