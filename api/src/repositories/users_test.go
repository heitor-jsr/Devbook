package repositories

import (
	"api/src/models"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/testcontainers/testcontainers-go"
)

type UserRepositorySuite struct {
	suite.Suite
	db        *sql.DB
	userRepo  *Usuarios
	container testcontainers.Container
	dsn       string
	tx        *sql.Tx
}

func (suite *UserRepositorySuite) SetupSuite() {
	ctx := context.Background()
	container, dsn, teardown := StartMySQLContainer(ctx, suite.T())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = db
	suite.userRepo = NewUsersRepository(db)
	suite.container = container
	suite.dsn = dsn

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

func (suite *UserRepositorySuite) SetupTest() {
	fmt.Println("SetUpTest...")

	erro := suite.db.Ping()

	if erro != nil {
		ctx := context.Background()
		container, dsn, _ := StartMySQLContainer(ctx, suite.T())

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			suite.T().Fatal(err)
		}

		suite.db = db
		suite.userRepo = NewUsersRepository(db)
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

func (suite *UserRepositorySuite) TestCreate() {
	t := suite.T()
	t.Run("Success on creating one first user", func(t *testing.T) {
		user := models.User{
			Nome:  "User Name",
			Nick:  "username",
			Email: "username@example.com",
			Senha: "password",
		}

		userID, err := suite.userRepo.Create(user)

		assert.NotNil(t, userID)
		assert.NoError(t, err)
		assert.IsType(t, uint64(3), userID)
		assert.Equal(t, uint64(3), userID)
	})

	t.Run("Fails when user.Nick already exists", func(t *testing.T) {
		user := models.User{
			Nome:  "John Does",
			Nick:  "johndoe",
			Email: "johndoe123@example.com",
			Senha: "password",
		}

		userID, err := suite.userRepo.Create(user)

		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), userID)
		assert.EqualError(t, err, "Error 1062 (23000): Duplicate entry 'johndoe' for key 'usuarios.nick'")
	})

	t.Run("Fails when user.Email already exists", func(t *testing.T) {
		user := models.User{
			Nome:  "John Does",
			Nick:  "johndoesss",
			Email: "johndoe@example.com",
			Senha: "password",
		}

		userID, err := suite.userRepo.Create(user)

		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), userID)
		assert.EqualError(t, err, "Error 1062 (23000): Duplicate entry 'johndoe@example.com' for key 'usuarios.email'")
	})

	t.Run("Fails when user.Nome is empty", func(t *testing.T) {
		user := models.User{
			Nick:  "johndoes",
			Email: "johndoe123@example.com",
			Senha: "password",
		}

		userID, err := suite.userRepo.Create(user)

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

		userID, err := suite.userRepo.Create(user)

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

		userID, err := suite.userRepo.Create(user)

		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), userID)
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()

		user := models.User{
			Nome:  "User Name",
			Nick:  "username",
			Email: "username@example.com",
			Senha: "password",
		}

		userID, err := suite.userRepo.Create(user)

		assert.Zero(t, userID)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func (suite *UserRepositorySuite) TestGetAll() {
	t := suite.T()
	t.Run("Success without name filter", func(t *testing.T) {
		users, err := suite.userRepo.GetAll("")

		assert.NotNil(t, users)
		assert.NoError(t, err)
		assert.IsType(t, []models.User{}, users)
		assert.Len(t, users, 2)
	})

	t.Run("Success with name filter", func(t *testing.T) {
		users, err := suite.userRepo.GetAll("The Second")

		assert.NotNil(t, users)
		assert.NoError(t, err)
		assert.IsType(t, []models.User{}, users)
		assert.Len(t, users, 1)
		assert.Contains(t, users, models.User{
			Id:       2,
			Nome:     "John Doe The Second",
			Nick:     "johndoe2",
			Email:    "johndoe2@example.com",
			CriadoEm: time.Time{},
		})
	})

	t.Run("Fails when the filter doesn't match any user", func(t *testing.T) {
		users, err := suite.userRepo.GetAll("zzz")

		assert.Nil(t, users)
		assert.NoError(t, err)
		assert.IsType(t, []models.User{}, users)
		assert.Len(t, users, 0)
	})
}

func (suite *UserRepositorySuite) TestGetById() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		user, err := suite.userRepo.GetById(2)

		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.IsType(t, models.User{}, user)
		assert.Exactly(t, user, models.User{
			Id:       2,
			Nome:     "John Doe The Second",
			Nick:     "johndoe2",
			Email:    "johndoe2@example.com",
			CriadoEm: time.Time{},
		})
	})

	t.Run("Fails", func(t *testing.T) {
		user, err := suite.userRepo.GetById(3)

		assert.NotNil(t, err)
		assert.Error(t, err, "no user found with ID: 3")
		assert.IsType(t, models.User{}, user)
		assert.Equal(t, fmt.Errorf("no user found with ID: %d", 3), err)
	})
}

func (suite *UserRepositorySuite) TestUpdate() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		user := models.User{
			Nome:  "John Doe The First",
			Nick:  "johndoefirst",
			Email: "johndoe@example.com",
		}

		err := suite.userRepo.Update(1, user)

		assert.Nil(t, err)
		assert.NoError(t, err)

		userUpdated, err := suite.userRepo.GetById(1)

		assert.IsType(t, models.User{}, userUpdated)
		assert.Exactly(t, userUpdated, models.User{
			Id:       1,
			Nome:     "John Doe The First",
			Nick:     "johndoefirst",
			Email:    "johndoe@example.com",
			CriadoEm: time.Time{},
		})
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()

		user := models.User{
			Nome:  "John Doe The First",
			Nick:  "johndoefirst",
			Email: "johndoe@example.com",
		}

		err := suite.userRepo.Update(2, user)

		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func (suite *UserRepositorySuite) TestGetByEmail() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		user, err := suite.userRepo.GetByEmail("johndoe2@example.com")
		assert.NoError(t, err)
		assert.IsType(t, models.User{}, user)
		assert.Exactly(t, user, models.User{
			Id:    2,
			Senha: "password",
		})
	})
}

func (suite *UserRepositorySuite) TestDelete() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		err := suite.userRepo.Delete(4)

		assert.Nil(t, err)
		assert.NoError(t, err)
	})
}

func (suite *UserRepositorySuite) TestFollowUser() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		err := suite.userRepo.Follow(2, 1)

		user, err := suite.userRepo.GetFollowers(1)

		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.IsType(t, []models.User{}, user)
		assert.Contains(t, user, models.User{
			Id:       2,
			Nome:     "John Doe The Second",
			Nick:     "johndoe2",
			Email:    "johndoe2@example.com",
			CriadoEm: time.Time{},
		})
	})
}

func (suite *UserRepositorySuite) TestUnfollowUser() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		err := suite.userRepo.Follow(2, 1)

		err2 := suite.userRepo.Unfollow(2, 1)

		user, err := suite.userRepo.GetFollowers(1)

		assert.Nil(t, err2, err)
		assert.NoError(t, err2, err)
		assert.IsType(t, []models.User{}, user)
	})
}

func (suite *UserRepositorySuite) TestGetFollowers() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		err := suite.userRepo.Follow(2, 1)

		user, err := suite.userRepo.GetFollowers(1)

		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.IsType(t, []models.User{}, user)
		assert.Contains(t, user, models.User{
			Id:       2,
			Nome:     "John Doe The Second",
			Nick:     "johndoe2",
			Email:    "johndoe2@example.com",
			CriadoEm: time.Time{},
		})
	})

	t.Run("Returns nil if user has no followers", func(t *testing.T) {
		user, err := suite.userRepo.GetFollowers(3)

		assert.IsType(t, []models.User{}, user)
		assert.Empty(t, user)
		assert.Nil(t, err)
		assert.NoError(t, err)
	})
}

func (suite *UserRepositorySuite) TestGetFollowing() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		err := suite.userRepo.Follow(2, 1)

		user, err := suite.userRepo.GetFollowing(2)
		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.IsType(t, []models.User{}, user)
		assert.Contains(t, user, models.User{
			Id:       1,
			Nome:     "John Doe",
			Nick:     "johndoe",
			Email:    "johndoe@example.com",
			CriadoEm: time.Time{},
		})
	})

	t.Run("Returns nil if user does not follow anyone", func(t *testing.T) {
		user, err := suite.userRepo.GetFollowing(3)

		assert.IsType(t, []models.User{}, user)
		assert.Empty(t, user)
		assert.Nil(t, err)
		assert.NoError(t, err)
	})
}

func (suite *UserRepositorySuite) TestGetPasswordFromDb() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		password, err := suite.userRepo.GetPasswordFromDb(1)

		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, "password", password)
	})

	t.Run("Returns empty string if user does not exists", func(t *testing.T) {
		password, err := suite.userRepo.GetPasswordFromDb(999)

		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, "", password)
		assert.Empty(t, password)
	})
}

func (suite *UserRepositorySuite) TestChangePassword() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		err := suite.userRepo.ChangePassword(1, "passwordChanged")

		assert.Nil(t, err)
		assert.NoError(t, err)

		password, err := suite.userRepo.GetPasswordFromDb(1)

		assert.Equal(t, "passwordChanged", password)
	})

	t.Run("Returns empty string if user does not exists", func(t *testing.T) {
		password, err := suite.userRepo.GetPasswordFromDb(999)

		assert.Nil(t, err)
		assert.NoError(t, err)
		assert.Equal(t, "", password)
		assert.Empty(t, password)
	})
}

func (suite *UserRepositorySuite) SeedDatabase() error {
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

func (suite *UserRepositorySuite) CleanDatabase() error {
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

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
