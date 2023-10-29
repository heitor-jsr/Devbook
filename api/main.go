package main

import (
	"api/src/router"
	"fmt"
	"net/http"
	"log"
)

// é esse o arquivo que vai ser o responsavel por executar o projeto. é ele que chama os outros pacotes, como o config, router, etc.

func main() {
	fmt.Println("API is running on port 5000")

	r := router.Generator()

	log.Fatal(http.ListenAndServe(":5000", r))
}