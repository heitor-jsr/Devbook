package controllers

import (
	"api/src/auth"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// a lógica do login é a seguinte: recebemos uma req com email e senha, verificamos se o usuario existe, se existe, verificamos se a senha comparada com o hash que tá no db é a correta e, se tudo estiver correto, retornamos um token de acesso. com isso o usuário estará logado no sistema.
func Login(w http.ResponseWriter, r *http.Request) {
	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	// pelo fato de ser uma rota que utiliza o post, a rota de login vai precisar dar o unmarshal nos dados do corpo da req e jgoar na variavel auxiliar, como todas as outras.
	var user models.User
	if erro = json.Unmarshal(body, &user); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	// a unica diferença desse método post é que precisamos buscar o usuario no db pelo email e, se encontrar, comparar a senha fornecida pelo usuario com a senha hasheada no db, para, se estiver tudo ok, retornar um token de acesso.
	userFromDB , erro := repository.GetByEmail(user.Email)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// o primeiro parametro é a senha do banco com o hash, e a segunda é a senha que o usuario passou ao tentar fazer o login.
	if erro = security.VerifyPassword(userFromDB.Senha, user.Senha); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	token ,_ := auth.GenerateToken(userFromDB.Id)
	fmt.Println("token: ",token)
	// responses.JSON(w, http.StatusOK, userFromDB
}