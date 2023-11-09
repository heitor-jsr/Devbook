package repositories

import (
	"api/src/models"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestUserRepository(t *testing.T) {

	ctx := context.Background()

	sqlScriptPath := filepath.Join("..", "..", "sql", "sql.sql")

	mysqlContainer, err := mysql.RunContainer(ctx,
    testcontainers.WithImage("mysql:8"),
    mysql.WithDatabase("devbook"),
    mysql.WithUsername("golang"),
    mysql.WithPassword("golang"),
		mysql.WithScripts(sqlScriptPath))
		if err != nil {
				panic(err)
		}

	defer func() {
			if err := mysqlContainer.Terminate(ctx); err != nil {
					panic(err)
			}
	}()

	containerPort, err := mysqlContainer.MappedPort(ctx, "3306/tcp")
	if err != nil {
			t.Fatal(err)
	}

	dsn := fmt.Sprintf("golang:golang@tcp(127.0.0.1:%s)/devbook", containerPort.Port())	

  db, err := sql.Open("mysql", dsn)
	if err != nil {
			t.Fatal(err)
	}
	defer db.Close()

	userRepo := NewUsersRepository(db)
		
	t.Run("CreateUser", func(t *testing.T) {
    user := models.User{
        Nome:  "John Doe",
        Nick:  "johndoe",
        Email: "johndoe@example.com",
        Senha: "password",
    }

    userID, err := userRepo.Create(user)

    if err != nil {
        t.Errorf("Error creating user: %v", err)
    }

    if userID == 0 {
        t.Error("User ID should not be 0")
    }
	})

	t.Run("CreateTwoUsers", func(t *testing.T) {
    user := models.User{
        Nome:  "John Doe The Seccond",
        Nick:  "johndoe2",
        Email: "johndoe2@example.com",
        Senha: "password",
    }

    userID, err := userRepo.Create(user)

    if err != nil {
        t.Errorf("Error creating user: %v", err)
    }

    if userID == 1 || userID == 0 {
        t.Error("User ID should be 2")
    }
		fmt.Println(userID, "dois")
	})
	t.Run("CreateUser fails when Nick already exists", func(t *testing.T) {
    user := models.User{
        Nome:  "John Does",
        Nick:  "johndoe",
        Email: "johndoe123@example.com",
        Senha: "password",
    }

    userID, err := userRepo.Create(user)

    if err == nil {
        t.Error("Expected an error, but got nil")
    }

    if userID == 3 {
        t.Error("User ID should not be 3")
    }
	})

	t.Run("CreateUser fails when user.Nome is empty", func(t *testing.T) {
    user := models.User{
        Nick:  "johndoes",
        Email: "johndoe123@example.com",
        Senha: "password",
    }

    userID, err := userRepo.Create(user)

    if err == nil {
        t.Errorf("Expected an error, but got nil, userID: %d", userID)
    }
	})

}
