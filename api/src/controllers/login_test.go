package controllers_test

import (
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
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type TestControllerSuite struct {
	suite.Suite
	db        *sql.DB
	userRepo  *repositories.Usuarios
	container testcontainers.Container
	dsn       string
	tx        *sql.Tx
}

func (suite *TestControllerSuite) SetupSuite() {
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

func (suite *TestControllerSuite) SetupTest() {
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

func (suite *TestControllerSuite) SeedDatabase() error {
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

func (suite *TestControllerSuite) CleanDatabase() error {
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

func (suite *TestControllerSuite) TestLogin() {
	t := suite.T()
	t.Run("Success on login", func(t *testing.T) {
	loginRequestBody := `{"email": "johndoe@example.com", "senha": "123456"}`

	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(loginRequestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	controllers.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var authData models.AuthenticationData
	err = json.NewDecoder(rr.Body).Decode(&authData)
	assert.NoError(t, err)

	assert.NotEmpty(t, authData.ID)
	assert.NotEmpty(t, authData.Token)
	})

	t.Run("Fail when request boyd is invalid", func(t *testing.T) {
		invalidRequestBody := "invalid_json"
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(invalidRequestBody))
		if err != nil {
			t.Fatal(err)
		}
	
		rr := httptest.NewRecorder()
		
		controllers.Login(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expectedResponseBody := "{\"erro\":\"invalid character 'i' looking for beginning of value\"}\n"
		assert.Equal(t, expectedResponseBody, rr.Body.String())
	
	})

	t.Run("Fail when user email is not found", func(t *testing.T) {
		loginRequestBody := `{"email": "johndoe123@example.com", "senha": "123456"}`
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(loginRequestBody))
		if err != nil {
			t.Fatal(err)
		}
	
		rr := httptest.NewRecorder()
		
		controllers.Login(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		expectedResponseBody := "{\"erro\":\"user with email johndoe123@example.com not found\"}\n"

		assert.Equal(t, expectedResponseBody, rr.Body.String())
	})

	t.Run("Fail when user password is incorrect", func(t *testing.T) {
		loginRequestBody := `{"email": "johndoe@example.com", "senha": "12345"}`
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(loginRequestBody))
		if err != nil {
			t.Fatal(err)
		}
	
		rr := httptest.NewRecorder()
		
		controllers.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		expectedResponseBody := "{\"erro\":\"crypto/bcrypt: hashedPassword is not the hash of the given password\"}\n"

		assert.Equal(t, expectedResponseBody, rr.Body.String())
	})
}

func TestLoginController(t *testing.T) {
	suite.Run(t, new(TestControllerSuite))
}