package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// para carregar essas variaveis de ambiente para dentro do go, nós precisamos de um pacote que se chama dotenv. ele é usado apenas para ler esses arquivos .env.

var (
	ConnStr = ""
	Port = 0
	SecretKey []byte
)

// vai carregar as variaveis de ambiente dentro da nossa aplicação. não recebe e nem retorna parametros. ela só vai alterar as variaveis de ambiente criadas antes dela, que vão estar disponíveis para a api toda. basicamente, a função abaixo vai só jogar valores dentro das variaveis de ambiente.
func Load() {
	var err error

	// abaixo usamos uma função do pacote dotenv que vai ler as informações do arquivo .env e disponibilizar elas no nosso ambiente de desenvolvimento. se ocorrer algum erro, a execução do programa deve ser encerrada, e a API nem deve subir. por isso do log.fatal. quem chama isso é o package main. se tudo der certo, as informações do .env já estão disponíveis para serem usadas e podem ser inseridas nas variaveis inicializadas acima.
	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}

	ConnStr = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}