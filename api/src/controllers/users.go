package controllers

import (
	"api/src/auth"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type UserController struct {
	db *sql.DB
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{db}
}

// o controller é o responsável por lidar com as requisições http e criar as respostas para os usuários.
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if erro := json.NewDecoder(r.Body).Decode(&user); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// antes mesmo da conexão com o db, devemos verificar se os dados recebidos do usuario sao validos, de acordo com os métodos do struct de users.
	if erro := user.Prepare("create"); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// não é responsabilidade do controller inserir dados no banco, ou manipular ele. como já foi dito, essa responsabilidade é atribuída ao repository. com ossio, o controller fica responsavel por processar a request, e solicitar a função adeequada do repository para manipular o db.

	// o fluxo é o seguinte: vc abre a conexão com o db, repassa essa conexao para o repositorio, e o metodo do repositorio responsavel por criar um novo usuario é chamado aqui, enviando os dados da request para que o repositorio venha a interagir com o db.
	repository := repositories.NewUsersRepository(uc.db)

	id, erro := repository.Create(user)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, id)
}

// busca todos os usuários de acordo com um filtro específico.
func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	// como sabemos, para capturar o parametro de busca de uma rota, precisamos primeiro capturar ele de forma dinâmica. parra isso, usamos o query string, que em go é feito da maneira abaixo.
	nameOrNick := strings.ToLower(r.URL.Query().Get("users"))

	repository := repositories.NewUsersRepository(uc.db)

	users, erro := repository.GetAll(nameOrNick)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func (uc UserController) GetUSerById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}
	repository := repositories.NewUsersRepository(uc.db)

	user, erro := repository.GetById(id)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	userIdInToken, erro := auth.ExtractUserId(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// vai bloquear que o usuário que está autorizado com o token gerado para fazer a operação tente realizar uma alteração em usuario que não é seu.
	if userIdInToken != id {
		responses.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário que não é o seu."))
		return
	}

	var user models.User
	if erro := json.NewDecoder(r.Body).Decode(&user); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = user.Prepare("update"); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repository := repositories.NewUsersRepository(uc.db)

	if erro = repository.Update(id, user); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// a função JSON criada vai receber, dentre outras coisas, os dados a serem enviados como resposta. porem, quando usamos um status code de noContent, não podemos passar nem o nil para ele. neses casos, precisamos fazer uma leve alteração no código do JSON, conforme é possivel ver nele, colocando o Encode() dentro de um bloco if data == nil {}.
	responses.JSON(w, http.StatusNoContent, nil)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	userIdInToken, erro := auth.ExtractUserId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if userIdInToken != id {
		responses.Erro(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não é o seu."))
		return
	}

	repository := repositories.NewUsersRepository(uc.db)

	if erro = repository.Delete(id); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (uc UserController) FollowUser(w http.ResponseWriter, r *http.Request) {
	// nessa rota nós devemos ter o seguinte cuidado: precisamos tanto do Id do seguidor, quanto do seguido. para isso, a lógica é simples: vamos pegar o userId do seguidor do token, e o userId do seguido do parametro da rota.
	followerId, erro := auth.ExtractUserId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	followedId, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followedId == followerId {
		responses.Erro(w, http.StatusForbidden, errors.New("Não pode seguir a si mesmo."))
		return
	}

	repository := repositories.NewUsersRepository(uc.db)

	if erro = repository.Follow(followerId, followedId); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (uc UserController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	// aqui precisamos dos dois ids, como na rota anterior. para isso, aplicamos a mesma lógica. relemebre que o id do token é o usuario que realiza para realizar a operação. é, portanto, o seguidor. o id do seguido é o que o usuario passa na rota.
	followerId, erro := auth.ExtractUserId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	followedId, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if followedId == followerId {
		responses.Erro(w, http.StatusForbidden, errors.New("Não pode deixar de seguir a si mesmo."))
		return
	}

	repository := repositories.NewUsersRepository(uc.db)

	if erro = repository.Unfollow(followerId, followedId); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (uc UserController) GetFollowers(w http.ResponseWriter, r *http.Request) {
	// não vamos pegar o id do token pq nós queremos permitir que um usuario busque os seguidores de outro usuario. por isso, vamos utilizar o id que vem na rota.
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repository := repositories.NewUsersRepository(uc.db)

	followers, erro := repository.GetFollowers(id)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func (uc UserController) GetFollowings(w http.ResponseWriter, r *http.Request) {
	// não vamos pegar o id do token pq você queremos permitir que um usuario busque os seguidores de outro usuario. por isso, vamos utilizar o id que vem na rota.
	params := mux.Vars(r)
	id, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repository := repositories.NewUsersRepository(uc.db)

	following, erro := repository.GetFollowing(id)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, following)
}

func (uc UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userIdFromToken, erro := auth.ExtractUserId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	userIdFromParams, erro := strconv.ParseUint(params["userId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if userIdFromToken != userIdFromParams {
		responses.Erro(w, http.StatusForbidden, errors.New("Não é possível alterar a senha de um usuário que não é o seu."))
		return
	}

	var password models.Password
	if erro := json.NewDecoder(r.Body).Decode(&password); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repository := repositories.NewUsersRepository(uc.db)


	// antes de alterar a senha, é indispensável que seja verificada se a senha atual que está sendo passada bate com a senha salva no db.
	passwordFromDb, erro := repository.GetPasswordFromDb(userIdFromToken)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// a senha que nos é retornada do banco está com hash. por isso, precisamos comparar a senha hasheada com a senha que o usuario passou na requisição e ver se seus valores batem.
	if erro := security.VerifyPassword(passwordFromDb, password.Current); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, errors.New("Senha atual inválida."))
		return
	}

	// se a senha do db e a atual q é passada na req forem iguais, precisamos, antes de salvar a nova senha no db, passar um hash pra ela.
	hashedPassowrd, erro := security.Hash(password.New)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// salva a nova senha no db
	if erro = repository.ChangePassword(userIdFromToken, string(hashedPassowrd)); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
