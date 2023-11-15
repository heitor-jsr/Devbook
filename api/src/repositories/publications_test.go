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

	// "time"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/testcontainers/testcontainers-go"
)

type PublicationsRepositorySuite struct {
	suite.Suite
	db        *sql.DB
	publiRepo  *Publications
	container testcontainers.Container
	dsn       string
	tx        *sql.Tx
}

func (suite *PublicationsRepositorySuite) SetupSuite() {
	ctx := context.Background()
	container, dsn, teardown := StartMySQLContainer(ctx, suite.T())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = db
	suite.publiRepo = NewPublicationRepository(db)
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

func (suite *PublicationsRepositorySuite) SetupTest() {
	erro := suite.db.Ping()

	if erro != nil {
		ctx := context.Background()
		container, dsn, _ := StartMySQLContainer(ctx, suite.T())

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			suite.T().Fatal(err)
		}

		suite.db = db
		suite.publiRepo = NewPublicationRepository(db)
		suite.container = container
		suite.dsn = dsn
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

func (suite *PublicationsRepositorySuite) TestCreate() {
	t := suite.T()
	t.Run("Success on creating one publication", func(t *testing.T) {
		publication := models.Publication{
			Title:   "Title",
			Content: "Content",
			AuthorId: 1,
		}

		publiCreated, err := suite.publiRepo.Create(publication)

		assert.NotNil(t, publiCreated)
		assert.NoError(t, err)
		assert.IsType(t, uint64(3), publiCreated)
		assert.Equal(t, uint64(3), publiCreated)
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()

		publication := models.Publication{
			Title:   "Title",
			Content: "Content",
			AuthorId: 1,
		}

		publiCreated, err := suite.publiRepo.Create(publication)

		assert.Zero(t, publiCreated)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

// func (suite *PublicationsRepositorySuite) TestGetAll() {
// 	t := suite.T()
// 	t.Run("Success without name filter", func(t *testing.T) {
// 		users, err := suite.publiRepo.GetAll("")

// 		assert.NotNil(t, users)
// 		assert.NoError(t, err)
// 		assert.IsType(t, []models.User{}, users)
// 		assert.Len(t, users, 2)
// 	})

// 	t.Run("Success with name filter", func(t *testing.T) {
// 		users, err := suite.publiRepo.GetAll("The Second")

// 		assert.NotNil(t, users)
// 		assert.NoError(t, err)
// 		assert.IsType(t, []models.User{}, users)
// 		assert.Len(t, users, 1)
// 		assert.Contains(t, users, models.User{
// 			Id:       2,
// 			Nome:     "John Doe The Second",
// 			Nick:     "johndoe2",
// 			Email:    "johndoe2@example.com",
// 			CriadoEm: time.Time{},
// 		})
// 	})

// 	t.Run("Fails when the filter doesn't match any user", func(t *testing.T) {
// 		users, err := suite.publiRepo.GetAll("zzz")

// 		assert.Nil(t, users)
// 		assert.NoError(t, err)
// 		assert.IsType(t, []models.User{}, users)
// 		assert.Len(t, users, 0)
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestGetById() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		user, err := suite.publiRepo.GetById(2)

// 		assert.NotNil(t, user)
// 		assert.NoError(t, err)
// 		assert.IsType(t, models.User{}, user)
// 		assert.Exactly(t, user, models.User{
// 			Id:       2,
// 			Nome:     "John Doe The Second",
// 			Nick:     "johndoe2",
// 			Email:    "johndoe2@example.com",
// 			CriadoEm: time.Time{},
// 		})
// 	})

// 	t.Run("Fails", func(t *testing.T) {
// 		user, err := suite.publiRepo.GetById(3)

// 		assert.NotNil(t, err)
// 		assert.Error(t, err, "no user found with ID: 3")
// 		assert.IsType(t, models.User{}, user)
// 		assert.Equal(t, fmt.Errorf("no user found with ID: %d", 3), err)
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestUpdate() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		user := models.User{
// 			Nome:  "John Doe The First",
// 			Nick:  "johndoefirst",
// 			Email: "johndoe@example.com",
// 		}

// 		err := suite.publiRepo.Update(1, user)

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)

// 		userUpdated, err := suite.publiRepo.GetById(1)

// 		assert.IsType(t, models.User{}, userUpdated)
// 		assert.Exactly(t, userUpdated, models.User{
// 			Id:       1,
// 			Nome:     "John Doe The First",
// 			Nick:     "johndoefirst",
// 			Email:    "johndoe@example.com",
// 			CriadoEm: time.Time{},
// 		})
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestGetByEmail() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		user, err := suite.publiRepo.GetByEmail("johndoe2@example.com")
// 		assert.NoError(t, err)
// 		assert.IsType(t, models.User{}, user)
// 		assert.Exactly(t, user, models.User{
// 			Id:    2,
// 			Senha: "password",
// 		})
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestDelete() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		err := suite.publiRepo.Delete(4)

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestFollowUser() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		err := suite.publiRepo.Follow(2, 1)

// 		user, err := suite.publiRepo.GetFollowers(1)

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 		assert.IsType(t, []models.User{}, user)
// 		assert.Contains(t, user, models.User{
// 			Id:       2,
// 			Nome:     "John Doe The Second",
// 			Nick:     "johndoe2",
// 			Email:    "johndoe2@example.com",
// 			CriadoEm: time.Time{},
// 		})
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestUnfollowUser() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		err := suite.publiRepo.Follow(2, 1)

// 		err2 := suite.publiRepo.Unfollow(2, 1)

// 		user, err := suite.publiRepo.GetFollowers(1)

// 		assert.Nil(t, err2, err)
// 		assert.NoError(t, err2, err)
// 		assert.IsType(t, []models.User{}, user)
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestGetFollowers() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		err := suite.publiRepo.Follow(2, 1)

// 		user, err := suite.publiRepo.GetFollowers(1)

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 		assert.IsType(t, []models.User{}, user)
// 		assert.Contains(t, user, models.User{
// 			Id:       2,
// 			Nome:     "John Doe The Second",
// 			Nick:     "johndoe2",
// 			Email:    "johndoe2@example.com",
// 			CriadoEm: time.Time{},
// 		})
// 	})

// 	t.Run("Returns nil if user has no followers", func(t *testing.T) {
// 		user, err := suite.publiRepo.GetFollowers(2)

// 		assert.IsType(t, []models.User{}, user)
// 		assert.Empty(t, user)
// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestGetFollowing() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		err := suite.publiRepo.Follow(2, 1)

// 		user, err := suite.publiRepo.GetFollowing(2)
// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 		assert.IsType(t, []models.User{}, user)
// 		assert.Contains(t, user, models.User{
// 			Id:       1,
// 			Nome:     "John Doe",
// 			Nick:     "johndoe",
// 			Email:    "johndoe@example.com",
// 			CriadoEm: time.Time{},
// 		})
// 	})

// 	t.Run("Returns nil if user does not follow anyone", func(t *testing.T) {
// 		user, err := suite.publiRepo.GetFollowing(1)

// 		assert.IsType(t, []models.User{}, user)
// 		assert.Empty(t, user)
// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestGetPasswordFromDb() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		password, err := suite.publiRepo.GetPasswordFromDb(1)

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "password", password)
// 	})

// 	t.Run("Returns empty string if user does not exists", func(t *testing.T) {
// 		password, err := suite.publiRepo.GetPasswordFromDb(999)

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "", password)
// 		assert.Empty(t, password)
// 	})
// }

// func (suite *PublicationsRepositorySuite) TestChangePassword() {
// 	t := suite.T()
// 	t.Run("Success", func(t *testing.T) {
// 		err := suite.publiRepo.ChangePassword(1, "passwordChanged")

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)

// 		password, err := suite.publiRepo.GetPasswordFromDb(1)

// 		assert.Equal(t, "passwordChanged", password)
// 	})

// 	t.Run("Returns empty string if user does not exists", func(t *testing.T) {
// 		password, err := suite.publiRepo.GetPasswordFromDb(999)

// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "", password)
// 		assert.Empty(t, password)
// 	})
// }


func (suite *PublicationsRepositorySuite) SeedDatabase() error {
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

func (suite *PublicationsRepositorySuite) CleanDatabase() error {
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

func TestPublicationsRepositorySuite(t *testing.T) {
	suite.Run(t, new(PublicationsRepositorySuite))
}
