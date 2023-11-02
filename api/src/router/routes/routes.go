// pacote que vai representar todas as rotas que teremos dentro do projeto. isso vai ser feito por meio de um struct.

package routes

import (
	"net/http"
	"github.com/gorilla/mux"
)

// Representa a estrutura de todas as rotas que serão implementadas pela API. todas as rotas do sistema vão ser criadas em cima desse struct.
type Route struct {
	URI string
	Method string
	Func func(w http.ResponseWriter, r *http.Request)
	RequireAuth bool
}

// vai configurar todas as rotas que temos e retornar elas para serem utilizadas no package router. ele vai receber um router como parametro, por meio de um ponteiro que aponta para o router que criamos e retornar esse mesmo ponteiro, indicando o local em memoria que ele está, mas agora com as rotas configuradas e prontas para uso.
func Config(r *mux.Router) *mux.Router{
	routes := usersRoutes
	routes = append(routes, loginRoute)
	// para lidar com o slice de rotas que criamos no package de routes, precisa de um for para cada item do slice. com isso, vamos percorrer cada rota, dando o handleFunc em cima das funções responsáveis por cada uma das rotas. esse é o motivo de termos criado um struct e um slice dele para armazenar as rotas.

	for _, route := range routes {
		r.HandleFunc(route.URI, route.Func).Methods(route.Method)
	}

	// o nosso r (router) entrou na função sem nenhuma rota, e é retornado com todas as rotas configuradas.
	return r
}