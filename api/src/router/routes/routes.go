// pacote que vai representar todas as rotas que teremos dentro do projeto. isso vai ser feito por meio de um struct.

package routes

import (
	"api/src/middlewares"
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
	routes = append(routes, publicationsRoutes...)
	// para lidar com o slice de rotas que criamos no package de routes, precisa de um for para cada item do slice. com isso, vamos percorrer cada rota, dando o handleFunc em cima das funções responsáveis por cada uma das rotas. esse é o motivo de termos criado um struct e um slice dele para armazenar as rotas.

	for _, route := range routes {
		// se a rota em questão requerer a autenticação, vai ser chamado o middlewares. se a rota em questão não requer a autenticação, vai ser chamado o handleFunc.
		// a primeira coisa que fazemos é passar uma função com a mesma assinatura que o middleware exige para ele ser chamado. ao ser chamada essa função que é passada como parametro para o middleware, ele se encarrega de autenticar o usuario e, se tudo estiver correto, chamar o metodo da rota e continuar a execução do código. senão, vai retornar um erro de autenticação.
		if route.RequireAuth {
			r.HandleFunc(route.URI, 
				middlewares.Logger(
					middlewares.Authentication(route.Func))).Methods(route.Method)
			continue
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Func)).Methods(route.Method)
		}
	}

	// o nosso r (router) entrou na função sem nenhuma rota, e é retornado com todas as rotas configuradas.
	return r
}