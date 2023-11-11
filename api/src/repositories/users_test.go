package repositories

import (
	"api/src/models"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	t.Run("Create", func(t *testing.T) {
		t.Run("Success on creating one first user", func(t *testing.T) {
			user := models.User{
				Nome:  "John Doe",
				Nick:  "johndoe",
				Email: "johndoe@example.com",
				Senha: "password",
			}

			userID, err := userRepo.Create(user)

			assert.NotNil(t, userID)
			assert.NoError(t, err)
			assert.IsType(t, uint64(1), userID)
			assert.Equal(t, uint64(1), userID)
		})

		t.Run("Success on creating one second user", func(t *testing.T) {
			user := models.User{
				Nome:  "John Doe The Seccond",
				Nick:  "johndoe2",
				Email: "johndoe2@example.com",
				Senha: "password",
			}

			userID, err := userRepo.Create(user)

			assert.NotNil(t, userID)
			assert.NoError(t, err)
			assert.IsType(t, uint64(2), userID)
			assert.Equal(t, uint64(2), userID)
		})
		t.Run("Fails when user.Nick already exists", func(t *testing.T) {
			user := models.User{
				Nome:  "John Does",
				Nick:  "johndoe",
				Email: "johndoe123@example.com",
				Senha: "password",
			}

			userID, err := userRepo.Create(user)

            assert.NotNil(t, err)
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userID)

		})

		t.Run("Fails when user.Nome is empty", func(t *testing.T) {
			user := models.User{
				Nick:  "johndoes",
				Email: "johndoe123@example.com",
				Senha: "password",
			}

			userID, err := userRepo.Create(user)

            assert.NotNil(t, err)
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userID)

		})

		t.Run("Fails when user.Email is empty", func(t *testing.T) {
			user := models.User{
				Nome:  "Jhen Doe",
				Nick:  "jhendoe",
				Senha: "password",
			}

			userID, err := userRepo.Create(user)

            assert.NotNil(t, err)
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userID)

		})

		t.Run("Fails when user.Senha is empty", func(t *testing.T) {
			user := models.User{
				Nome:  "Jhen Doe",
				Nick:  "jhendoe",
				Email: "jhen123@example.com",
			}

			userID, err := userRepo.Create(user)

            assert.NotNil(t, err)
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userID)

		})
	})

  t.Run("GetAll", func(t *testing.T) {
		t.Run("Success without name filter", func(t *testing.T) {
			users, err := userRepo.GetAll("")

			assert.NotNil(t, users)
			assert.NoError(t, err)
			assert.IsType(t, []models.User{}, users)
			assert.Len(t, users, 2)

		})

		t.Run("Success with name filter", func(t *testing.T) {
			users, err := userRepo.GetAll("The Seccond")

			assert.NotNil(t, users)
			assert.NoError(t, err)
			assert.IsType(t, []models.User{}, users)
			assert.Len(t, users, 1)
			assert.Contains(t, users, models.User{
					Id:    2,
					Nome:  "John Doe The Seccond",
					Nick:  "johndoe2",
					Email: "johndoe2@example.com",
					CriadoEm: time.Time{},
			})
		})
		t.Run("Fails when the filter doesn't match any user", func(t *testing.T) {
			users, err := userRepo.GetAll("zzz")

      assert.Nil(t, users)
			assert.NoError(t, err)
			assert.IsType(t, []models.User{}, users)
			assert.Len(t, users, 0)

		})
	})	
		t.Run("GetById", func(t *testing.T) {
			t.Run("Success", func(t *testing.T) {
				user, err := userRepo.GetById(2)
	
				assert.NotNil(t, user)
				assert.NoError(t, err)
				assert.IsType(t, models.User{}, user)
				assert.Exactly(t, user, models.User{
					Id:    2,
					Nome:  "John Doe The Seccond",
					Nick:  "johndoe2",
					Email: "johndoe2@example.com",
					CriadoEm: time.Time{},
			})
		})
	
			t.Run("Fails", func(t *testing.T) {
				user, err := userRepo.GetById(3)
	
				assert.NotNil(t, err)
				assert.Error(t, err, "no user found with ID: 3")
				assert.IsType(t, models.User{}, user)
				assert.Equal(t, fmt.Errorf("no user found with ID: %d", 3), err)
			})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			user := models.User{
				Nome:  "John Doe The First",
				Nick:  "johndoefirst",
				Email: "johndoe@example.com",
			}

			err := userRepo.Update(1, user)

			assert.Nil(t, err)
			assert.NoError(t, err)

			userUpdated, err := userRepo.GetById(1)

			assert.IsType(t, models.User{}, userUpdated)
			assert.Exactly(t, userUpdated, models.User{
				Id:    1,
				Nome:  "John Doe The First",
				Nick:  "johndoefirst",
				Email: "johndoe@example.com",
				CriadoEm: time.Time{},
			})
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err := userRepo.Delete(1)

			assert.Nil(t, err)
			assert.NoError(t, err)
		})
	})
}
