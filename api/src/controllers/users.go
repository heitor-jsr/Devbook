package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// o controller é o responsável por lidar com as requisições http e criar as respostas para os usuários.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, erro := ioutil.ReadAll(r.Body)
	
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.User

	if erro = json.Unmarshal(reqBody, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// antes mesmo da conexão com o db, devemos verificar se os dados recebidos do usuario sao validos, de acordo com os métodos do struct de users.
	if erro = usuario.Prepare(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// não é responsabilidade do controller inserir dados no banco, ou manipular ele. como já foi dito, essa responsabilidade é atribuída ao repository. com ossio, o controller fica responsavel por processar a request, e solicitar a função adeequada do repository para manipular o db.
	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// o fluxo é o seguinte: vc abre a conexão com o db, repassa essa conexao para o repositorio, e o metodo do repositorio responsavel por criar um novo usuario é chamado aqui, enviando os dados da request para que o repositorio venha a interagir com o db.
	repository := repositories.NewUsersRepository(db)

	id, erro := repository.Create(usuario)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	
	responses.JSON(w, http.StatusCreated, id)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get all users"))
}

func GetUSerById(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get user by his id"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update user by his id"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete user"))
}