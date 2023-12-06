package controllers_test

import (
	"api/src/auth"
	"api/src/controllers"
	"api/src/models"
	"api/src/repositories"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type TestUserControllerSuite struct {
	suite.Suite
	db             *sql.DB
	userRepo       *repositories.Usuarios
	container      testcontainers.Container
	userController *controllers.UserController
	dsn            string
	tx             *sql.Tx
}

func (suite *TestUserControllerSuite) SetupSuite() {
	envPath := "/home/dornzak/devbook/api/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal(err)
	}

	mysqlUser := os.Getenv("DB_USER")
	mysqlPassword := os.Getenv("DB_PASSWORD")
	mysqlDbName := os.Getenv("DB_NAME")

	fmt.Println("DB_USER:", mysqlUser)
	fmt.Println("DB_PASSWORD:", mysqlPassword)
	fmt.Println("DB_NAME:", mysqlDbName)

	ctx := context.Background()
	container, dsn, teardown := repositories.StartMySQLContainer(ctx, suite.T())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = db
	suite.userRepo = repositories.NewUsersRepository(db)
	suite.container = container
	suite.dsn = dsn

	fmt.Println("DSN:", dsn)

	suite.T().Cleanup(func() {
		suite.T().Cleanup(func() {
			if err := suite.tx.Rollback(); err != nil {
				suite.T().Fatal(err)
			}

			_, err := teardown()
			if err != nil {
				suite.T().Fatal(err)
			}
		})
		if suite.db != nil {
			if err := suite.db.Close(); err != nil {
				suite.T().Errorf("Error closing database connection: %v", err)
			}
		}
	})

	suite.tx, err = suite.db.Begin()
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TestUserControllerSuite) SetupTest() {
	fmt.Println("SetUpTest...")

	erro := suite.db.Ping()

	if erro != nil {
		ctx := context.Background()
		container, dsn, _ := repositories.StartMySQLContainer(ctx, suite.T())

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			suite.T().Fatal(err)
		}

		suite.db = db
		suite.userRepo = repositories.NewUsersRepository(db)
		suite.container = container
		suite.dsn = dsn

		err = suite.CleanDatabase()
		if err != nil {
			suite.T().Fatal(err)
		}

		err = suite.SeedDatabase()
		if err != nil {
			suite.T().Fatal(err)
		}
	} else {
		err := suite.CleanDatabase()
		if err != nil {
			suite.T().Fatal(err)
		}

		err = suite.SeedDatabase()
		if err != nil {
			suite.T().Fatal(err)
		}
	}
}

func (suite *TestUserControllerSuite) SeedDatabase() error {
	fmt.Println("Seeding database...")
	insertDataScriptPath := filepath.Join("..", "..", "sql", "insert_data.sql")

	insertDataScript, err := ioutil.ReadFile(insertDataScriptPath)
	if err != nil {
		suite.T().Errorf("Erro ao ler script de inserção de dados: %v", err)
		return err
	}

	queries := strings.Split(string(insertDataScript), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)

		if query != "" {
			_, err := suite.db.ExecContext(context.Background(), query)
			if err != nil {
				suite.T().Errorf("Erro ao executar a instrução SQL: %v", err)
				return err
			}
		}
	}

	return nil
}

func (suite *TestUserControllerSuite) CleanDatabase() error {
	_, err := suite.db.ExecContext(context.Background(), "DELETE FROM usuarios")
	if err != nil {
		return err
	}

	_, err = suite.db.ExecContext(context.Background(), "ALTER TABLE usuarios AUTO_INCREMENT = 1")
	if err != nil {
		return err
	}

	return nil
}

func (suite *TestUserControllerSuite) TestCreate() {
	t := suite.T()
	t.Run("Success on creating user", func(t *testing.T) {
		createUserRequestBody := `{"Nome": "John Doe Forth", "Email": "john.doe.forth@example.com", "Senha": "strongPassword123", "Nick": "john_doe_forth"}`

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		var uc = controllers.NewUserController(suite.db)

		uc.CreateUser(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var lasInsertId uint64

		err = json.NewDecoder(rr.Body).Decode(&lasInsertId)
		assert.NoError(t, err)
		assert.NotEmpty(t, lasInsertId)
		assert.Equal(t, uint64(3), lasInsertId)
	})

	t.Run("Fail when name is missing", func(t *testing.T) {
		createUserRequestBody := `{"Email": "john.doe.third@example.com", "Senha": "strongPassword123", "Nick": "third"}`

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		var uc = controllers.NewUserController(suite.db)

		uc.CreateUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "o campo nome é obrigatorio",
		}

		expectedResponseMessage := "o campo nome é obrigatorio"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)
		assert.NotNil(t, rr.Body.String())

		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)

	})

	t.Run("Fail when email is missing", func(t *testing.T) {
		createUserRequestBody := `{"Nome": "John Doe Fifth", "Senha": "strongPassword123", "Nick": "john_doe_forth"}`

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		var uc = controllers.NewUserController(suite.db)

		uc.CreateUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "o campo email é obrigatorio",
		}

		expectedResponseMessage := "o campo email é obrigatorio"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)
		assert.NotNil(t, rr.Body.String())

		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)

	})

	t.Run("Fail when nick is missing", func(t *testing.T) {
		createUserRequestBody := `{"Nome": "John Doe Fifth", "Email": "john.doe.fifth@example.com", "Senha": "strongPassword123"}`

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		var uc = controllers.NewUserController(suite.db)

		uc.CreateUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "o campo nick é obrigatorio",
		}

		expectedResponseMessage := "o campo nick é obrigatorio"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)
		assert.NotNil(t, rr.Body.String())

		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)

	})

	t.Run("Fail when password is missing", func(t *testing.T) {
		createUserRequestBody := `{"Nome": "John Doe Fifth", "Email": "john.doe.fifth@example.com", "Nick": "john_doe_fifth"}`

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		var uc = controllers.NewUserController(suite.db)

		uc.CreateUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "o campo senha é obrigatorio",
		}

		expectedResponseMessage := "o campo senha é obrigatorio"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)
		assert.NotNil(t, rr.Body.String())

		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)

	})
}

func (suite *TestUserControllerSuite) TestGetUsers() {
	t := suite.T()
	t.Run("Success on getting user by name", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

		req, err := http.NewRequest("GET", "/users?name=John", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		uc.GetUsers(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var users []models.User
		err = json.NewDecoder(rr.Body).Decode(&users)
		if err != nil {
			t.Fatal(err)
		}

		var resposeExpected = []models.User{
			{
				Id:    1,
				Nome:  "John Doe",
				Email: "johndoe@example.com",
				Nick:  "johndoe",
				Senha: "",
			},
			{
				Id:    2,
				Nome:  "John Doe The Second",
				Email: "johndoe2@example.com",
				Nick:  "johndoe2",
				Senha: "",
			},
			{
				Id:    3,
				Nome:  "John Doe The Third",
				Email: "johndoe3@example.com",
				Nick:  "johndoe3",
				Senha: "",
			},
		}

		assert.Equal(t, 3, len(users))
		assert.Contains(t, users, resposeExpected[0])
		assert.Contains(t, users, resposeExpected[1])
		assert.Contains(t, users, resposeExpected[2])
		assert.ElementsMatch(t, resposeExpected, users)
		assert.EqualValues(t, resposeExpected, users)
	})

	t.Run("Success on getting user by nick", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

		req, err := http.NewRequest("GET", "/users?name=johndoe2", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		uc.GetUsers(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var users []models.User
		err = json.NewDecoder(rr.Body).Decode(&users)
		if err != nil {
			t.Fatal(err)
		}

		var resposeExpected = []models.User{
			{
				Id:    1,
				Nome:  "John Doe",
				Email: "johndoe@example.com",
				Nick:  "johndoe",
				Senha: "",
			},
			{
				Id:    2,
				Nome:  "John Doe The Second",
				Email: "johndoe2@example.com",
				Nick:  "johndoe2",
				Senha: "",
			},
			{
				Id:    3,
				Nome:  "John Doe The Third",
				Email: "johndoe3@example.com",
				Nick:  "johndoe3",
				Senha: "",
			},
		}

		assert.Equal(t, 3, len(users))
		assert.Contains(t, users, resposeExpected[0])
		assert.Contains(t, users, resposeExpected[1])
		assert.Contains(t, users, resposeExpected[2])
		assert.ElementsMatch(t, resposeExpected, users)
		assert.EqualValues(t, resposeExpected, users)
	})


	t.Run("Success on getting user when no parameters is given", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		uc.GetUsers(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var users []models.User
		err = json.NewDecoder(rr.Body).Decode(&users)
		if err != nil {
			t.Fatal(err)
		}

		var resposeExpected = []models.User{
			{
				Id:    1,
				Nome:  "John Doe",
				Email: "johndoe@example.com",
				Nick:  "johndoe",
				Senha: "",
			},
			{
				Id:    2,
				Nome:  "John Doe The Second",
				Email: "johndoe2@example.com",
				Nick:  "johndoe2",
				Senha: "",
			},			
			{
				Id:    3,
				Nome:  "John Doe The Third",
				Email: "johndoe3@example.com",
				Nick:  "johndoe3",
				Senha: "",
			},
		}

		assert.Equal(t, 3, len(users))
		assert.Contains(t, users, resposeExpected[0])
		assert.Contains(t, users, resposeExpected[1])
		assert.Contains(t, users, resposeExpected[2])
		assert.ElementsMatch(t, resposeExpected, users)
		assert.EqualValues(t, resposeExpected, users)
	})
}

func (suite *TestUserControllerSuite) TestGetUserById() {
	t := suite.T()
	t.Run("Success on getting user by id", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

		req, err := http.NewRequest("GET", "/users/1", nil)

		// Para testar requisições que utilizam o pacote mux e o roteador fornecido por ele, é necessário criar um roteador do mux no seu teste, adicionar a rota que está sendo testada e, em seguida, servir a requisição através desse roteador. Sem isso, o mux.Vars(r) não consegue extrair o parâmetro "userId" da URL, porque a requisição não passa por um roteador do mux.

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.GetUSerById)

		rr := httptest.NewRecorder()

    router.ServeHTTP(rr, req)

		uc.GetUSerById(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var users models.User
		err = json.NewDecoder(rr.Body).Decode(&users)
		if err != nil {
			t.Fatal(err)
		}

		var resposeExpected = models.User{
				Id:    1,
				Nome:  "John Doe",
				Email: "johndoe@example.com",
				Nick:  "johndoe",
				Senha: "",
		}

		assert.Equal(t, resposeExpected, users)
		assert.EqualValues(t, resposeExpected, users)
	})

	t.Run("Fail when string is not an integer", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

		req, err := http.NewRequest("GET", "/users/asdf", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.GetUSerById)

		rr := httptest.NewRecorder()

    router.ServeHTTP(rr, req)

		uc.GetUSerById(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "strconv.ParseUint: parsing \"asdf\": invalid syntax",
		}

		expectedResponseMessage := "strconv.ParseUint: parsing \"asdf\": invalid syntax"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)
		assert.NotNil(t, rr.Body.String())

		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)
	})
}

func (suite *TestUserControllerSuite) TestUpdateUser() {
	t := suite.T()
	t.Run("Success on updating user", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

    req, err := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"Nome": "João de tal", "Email": "joaodetal@example.com", "Senha": "strongPassword", "Nick": "jao_de_tal"}`))
		if err != nil {
			t.Fatal(err)
		}

		// Para testar requisições que utilizam o pacote mux e o roteador fornecido por ele, é necessário criar um roteador do mux no seu teste, adicionar a rota que está sendo testada e, em seguida, servir a requisição através desse roteador. Sem isso, o mux.Vars(r) não consegue extrair o parâmetro "userId" da URL, porque a requisição não passa por um roteador do mux.
		token, _ := auth.GenerateToken(1)

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.UpdateUser)

		rr := httptest.NewRecorder()
		
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

    router.ServeHTTP(rr, req)

		uc.UpdateUser(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Fail when the id from token is different from the id from URL", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

    req, err := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"Nome": "João de tal 2", "Email": "joaodetal2@example.com", "Senha": "strongPassword", "Nick": "jao_de_tal2"}`))
		if err != nil {
			t.Fatal(err)
		}

		tokenInvalid, _ := auth.GenerateToken(2)

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.UpdateUser)

		rr := httptest.NewRecorder()
		
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenInvalid))

    router.ServeHTTP(rr, req)

		uc.UpdateUser(rr, req)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "Não é possível atualizar um usuário que não é o seu.",
		}

		expectedResponseMessage := "Não é possível atualizar um usuário que não é o seu."

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.NotNil(t, responseReturned)
		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)

	})

	t.Run("Fail when extracting id from token", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

    req, err := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"Nome": "João de tal 2", "Email": "joaodetal2@example.com", "Senha": "strongPassword", "Nick": "jao_de_tal2"}`))
    if err != nil {
        t.Fatal(err)
    }

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "token_invalido"))

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.UpdateUser)

		rr := httptest.NewRecorder()
		
    router.ServeHTTP(rr, req)

		uc.UpdateUser(rr, req)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "token contains an invalid number of segments",
		}

		expectedResponseMessage := "token contains an invalid number of segments"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.NotNil(t, responseReturned)
		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)
	})
	t.Run("Fail when body is invalid", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

    req, err := http.NewRequest("PUT", "/users/1", strings.NewReader(`{"Nome": "João de tal 2", "Email": "joaodetal2@example.com", "Senha": "strongPassword"}`))
    if err != nil {
        t.Fatal(err)
    }

		token, _ := auth.GenerateToken(1)

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.UpdateUser)

		rr := httptest.NewRecorder()
		
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

    router.ServeHTTP(rr, req)

		uc.UpdateUser(rr, req)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "o campo nick é obrigatorio",
		}

		expectedResponseMessage := "o campo nick é obrigatorio"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotNil(t, responseReturned)
		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)
	})
}

func (suite *TestUserControllerSuite) TestDeleteUser() {
	t := suite.T()
	t.Run("Success on deleting user", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

    req, err := http.NewRequest("DELETE", "/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		token, _ := auth.GenerateToken(1)

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.DeleteUser)

		rr := httptest.NewRecorder()
		
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

    router.ServeHTTP(rr, req)

		uc.DeleteUser(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Fail when the id from token is different from the id from URL", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

    req, err := http.NewRequest("DELETE", "/users/2", nil)
		if err != nil {
			t.Fatal(err)
		}

		tokenInvalid, _ := auth.GenerateToken(1)

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.DeleteUser)

		rr := httptest.NewRecorder()
		
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenInvalid))

    router.ServeHTTP(rr, req)

		uc.DeleteUser(rr, req)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "Não é possível deletar um usuário que não é o seu.",
		}

		expectedResponseMessage := "Não é possível deletar um usuário que não é o seu."

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.NotNil(t, responseReturned)
		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)

	})

	t.Run("Fail when extracting id from token", func(t *testing.T) {
		var uc = controllers.NewUserController(suite.db)

    req, err := http.NewRequest("DELETE", "/users/2", nil)
    if err != nil {
        t.Fatal(err)
    }

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "token_invalido"))

		router := mux.NewRouter()
    router.HandleFunc("/users/{userId}", uc.DeleteUser)

		rr := httptest.NewRecorder()
		
    router.ServeHTTP(rr, req)

		uc.DeleteUser(rr, req)

		var responseReturned struct {
			Erro string
		}

		expectedResponseBody := struct {
			Erro string
		}{
			Erro: "token contains an invalid number of segments",
		}

		expectedResponseMessage := "token contains an invalid number of segments"

		err = json.NewDecoder(rr.Body).Decode(&responseReturned)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.NotNil(t, responseReturned)
		assert.Equal(t, expectedResponseMessage, responseReturned.Erro)
		assert.Equal(t, expectedResponseBody, responseReturned)
	})

}
func TestUserController(t *testing.T) {
	suite.Run(t, new(TestUserControllerSuite))
}
