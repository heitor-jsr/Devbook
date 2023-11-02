package main

import (
	"api/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

// é esse o arquivo que vai ser o responsavel por executar o projeto. é ele que chama os outros pacotes, como o config, router, etc.

// func init() {
// 	key := make([]byte, 64)
// 	// vai pegar o slice que só possui valores 0 e encher ele de números aleatórios cada vez que o programa for inicializado.
// 	if _, err := rand.Read(key); err != nil {
// 		log.Fatal(err)
// 	}

// 	// gerado esse nosso secret de maneira aleatória e segura, precisamos salva-lo no .env como uma string e não como um slice de byte.

// 	stringBase64 := base64.StdEncoding.EncodeToString(key)
	
// 	fmt.Println(stringBase64)
// }

func main() {
	// para carregar os arquivos .env, basta executarmos a função load. para ter certeza que tudo deu certo, precisamos apenas dar um log na porta. por fim, para dar o listen and serve na porta do .env, basta executar o seguinte comando:
	config.Load()

	fmt.Printf("API is running on port %d", config.Port)

	r := router.Generator()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}