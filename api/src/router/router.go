package router

import (
	"api/src/router/routes"

	"github.com/gorilla/mux"
)

// vai gerar e retornar um ponteiro para o router. com isso, retornamos um novo router, com as rotas já configuradas para serem utilizadas nos outros pacotes.
func Generator() *mux.Router{
	r := mux.NewRouter()
	
	// foi criado o nosso router sem nenhuma configuração de rota, e depois disso, retornamos (conforme a func Config do package routes) ele com a configuração das rotas que criamos no package de routes.
	return routes.Config(r)
}