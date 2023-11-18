package controllers

import (
	"api/config"
	"api/src/auth"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"fmt"

	// "io"
	"net/http"
	"strconv"
)

func init () {
	config.Load()
}

// a lógica do login é a seguinte: recebemos uma req com email e senha, verificamos se o usuario existe, se existe, verificamos se a senha comparada com o hash que tá no db é a correta e, se tudo estiver correto, retornamos um token de acesso. com isso o usuário estará logado no sistema.
func Login(w http.ResponseWriter, r *http.Request) {
	// pelo fato de ser uma rota que utiliza o post, a rota de login vai precisar dar o unmarshal nos dados do corpo da req e jgoar na variavel auxiliar, como todas as outras. alteramos o unmarshal para newdecoder, pois ele tem um gerenciamento de memoria mais eficaz.
	var user models.User
	if erro := json.NewDecoder(r.Body).Decode(&user); erro != nil {
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

	fmt.Println(userFromDB)

	// o primeiro parametro é a senha do banco com o hash, e a segunda é a senha que o usuario passou ao tentar fazer o login.
	if erro = security.VerifyPassword(userFromDB.Senha, user.Senha); erro != nil {
		fmt.Println("Error verifying password:", erro)
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	token, erro := auth.GenerateToken(userFromDB.Id)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	usuarioID := strconv.FormatUint(userFromDB.Id, 10)
	fmt.Println("Before responses.JSON")
	responses.JSON(w, http.StatusOK, models.AuthenticationData{ID: usuarioID, Token: token})
}