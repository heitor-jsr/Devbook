package main

import (
	"api/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

// é esse o arquivo que vai ser o responsavel por executar o projeto. é ele que chama os outros pacotes, como o config, router, etc.

func main() {
	// para carregar os arquivos .env, basta executarmos a função load. para ter certeza que tudo deu certo, precisamos apenas dar um log na porta. por fim, para dar o listen and serve na porta do .env, basta executar o seguinte comando:
	config.Load()

	fmt.Printf("API is running on port %d", config.Port)

	r := router.Generator()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}