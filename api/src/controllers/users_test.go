package controllers_test

import (
	"api/src/controllers"
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

	t.Run("Fail when nick is missing", func(t *testing.T) {
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

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		createUserRequestBody := `{"Nome": "John Doe Forth", "Email": "john.doe.forth@example.com", "Senha": "strongPassword123", "Nick": "john_doe_forth"}`

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserRequestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()


		var uc = controllers.NewUserController(suite.db)
		
		uc.CreateUser(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(TestUserControllerSuite))
}
