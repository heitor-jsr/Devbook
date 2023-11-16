package repositories

import (
	"api/src/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// "reflect"
	"strings"
	"testing"

	// "bou.ke/monkey"
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

type FakeResult struct {
	sql.Result
}

func (f *FakeResult) LastInsertId() (int64, error) {
	return 0, errors.New("Simulated LastInsertId error")
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
	fmt.Println("SetupTest")
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

	t.Run("Fail when try to create a publication without title", func(t *testing.T) {
		publication := models.Publication{
			Content: "Content",
			AuthorId: 1,
		}

		publiCreated, err := suite.publiRepo.Create(publication)

		assert.Zero(t, publiCreated)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error 3819 (HY000): Check constraint 'publicacoes_chk_1' is violated.")
	})

	t.Run("Fail when try to create a publication without content", func(t *testing.T) {
		publication := models.Publication{
			Title:   "Title",
			AuthorId: 1,
		}

		publiCreated, err := suite.publiRepo.Create(publication)

		assert.Zero(t, publiCreated)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error 3819 (HY000): Check constraint 'publicacoes_chk_2' is violated.")
	})

	t.Run("Fail when try to create a publication without author_id", func(t *testing.T) {
		publication := models.Publication{
			Title:   "Title",
			Content:   "Content",
		}

		publiCreated, err := suite.publiRepo.Create(publication)

		assert.Zero(t, publiCreated)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error 3819 (HY000): Check constraint 'publicacoes_chk_3' is violated.")
	})

	// t.Run("Fail when try to create a publication without title", func(t *testing.T) {
  //   monkey.PatchInstanceMethod(reflect.TypeOf(&FakeResult{}), "LastInsertId", func(_ *FakeResult) (int64, error) {
  //       return 0, errors.New("Simulated LastInsertId error")
  //   })

  //   publication := models.Publication{
  //       Title:    "Title",
  //       Content:  "Content",
  //       AuthorId: 1,
  //   }

  //   publiCreated, err := suite.publiRepo.Create(publication)

  //   assert.Zero(t, publiCreated)
  //   assert.NotNil(t, err)
  //   assert.Error(t, err)
  //   assert.Contains(t, err.Error(), "Simulated LastInsertId error")
  //   defer monkey.UnpatchAll()

	// })

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

func (suite *PublicationsRepositorySuite) TestGetAll() {
	t := suite.T()
	t.Run("Success on getting all publications", func(t *testing.T) {
		publication, err := suite.publiRepo.GetPublications(2)

		assert.NotNil(t, publication)
		assert.NoError(t, err)
		assert.IsType(t, []models.Publication{}, publication)
		assert.Len(t, publication, 2)
		assert.ElementsMatch(t, publication, []models.Publication{{
			Id:         2,
			Title:      "John Doe The Second",
			Content:    "John Doe The Second Content",
			AuthorId:   2,
			AuthorNick: "johndoe2",
			Likes:      0,
			CriadaEm:   publication[0].CriadaEm,
			}, {
			Id:         1,
			Title:      "John Doe Publication",
			Content:    "John Doe Publication Content",
			AuthorId:   1,
			AuthorNick: "johndoe",
			Likes:      0,
			CriadaEm:   publication[1].CriadaEm,
			},
		})	
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()

		publication, err := suite.publiRepo.GetPublications(2)

		assert.Zero(t, publication)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func (suite *PublicationsRepositorySuite) TestGetById() {
	t := suite.T()
	t.Run("Success on getting the publication by id", func(t *testing.T) {
		publication, err := suite.publiRepo.GetPublicationById(2)

		assert.NotNil(t, publication)
		assert.NoError(t, err)
		assert.IsType(t, models.Publication{}, publication)
		assert.ObjectsAreEqualValues(publication, models.Publication{
			Id:         2,
			Title:      "John Doe The Second",
			Content:    "John Doe The Second Content",
			AuthorId:   2,
			AuthorNick: "johndoe2",
			Likes:      0,
			CriadaEm:   publication.CriadaEm,
			},
		)	
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()

		publication, err := suite.publiRepo.GetPublicationById(2)

		assert.Zero(t, publication)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func (suite *PublicationsRepositorySuite) TestUpdate() {
	t := suite.T()
	t.Run("Success on updating the publication", func(t *testing.T) {
		publication := models.Publication{
			Title:   "John Doe The Second",
			Content: "John Doe The Second Content Updated",
			AuthorId: 2,
		}

		err := suite.publiRepo.Update(2, publication)

		publicationUpdated, err := suite.publiRepo.GetPublicationById(2)

		assert.NotNil(t, publicationUpdated)
		assert.NoError(t, err)
		assert.IsType(t, models.Publication{}, publicationUpdated)
		assert.ObjectsAreEqualValues(publicationUpdated, models.Publication{
			Id:         2,
			Title:      "John Doe The Second",
			Content:    "John Doe The Second Content Updated",
			AuthorId:   2,
			AuthorNick: "johndoe2",
			Likes:      0,
			CriadaEm:   publication.CriadaEm,
			},
		)	
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()

		publication := models.Publication{
			Title:   "John Doe The Second",
			Content: "John Doe The Second Content Updated",
			AuthorId: 2,
		}

		err := suite.publiRepo.Update(2, publication)

		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func (suite *PublicationsRepositorySuite) TestDelete() {
	t := suite.T()
	t.Run("Success", func(t *testing.T) {
		err := suite.publiRepo.Delete(1)

		assert.Nil(t, err)
		assert.NoError(t, err)
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()
		err := suite.publiRepo.Delete(1)


		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func (suite *PublicationsRepositorySuite) TestGetPublicationFromUser() {
	t := suite.T()
	t.Run("Success on getting the publication by user id", func(t *testing.T) {
		publication, err := suite.publiRepo.GetPublicationFromUser(2)

		assert.NotNil(t, publication)
		assert.NoError(t, err)
		assert.IsType(t, []models.Publication{}, publication)
		assert.ElementsMatch(t, publication, []models.Publication{{
			Id:         2,
			Title:      "John Doe The Second",
			Content:    "John Doe The Second Content",
			AuthorId:   2,
			AuthorNick: "johndoe2",
			Likes:      0,
			CriadaEm:   publication[0].CriadaEm,
			},
		})
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()

		publication, err := suite.publiRepo.GetPublicationFromUser(2)

		assert.Zero(t, publication)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

func (suite *PublicationsRepositorySuite) TestLikePublication() {
	t := suite.T()
	t.Run("Success on updating the publication", func(t *testing.T) {
		err := suite.publiRepo.LikePublication(2)

		publicationLiked, err := suite.publiRepo.GetPublicationById(2)

		assert.NotNil(t, publicationLiked)
		assert.NoError(t, err)
		assert.IsType(t, models.Publication{}, publicationLiked)
		assert.ObjectsAreEqualValues(publicationLiked, models.Publication{
			Id:         2,
			Title:      "John Doe The Second",
			Content:    "John Doe The Second Content Updated",
			AuthorId:   2,
			AuthorNick: "johndoe2",
			Likes:      1,
			CriadaEm:   publicationLiked.CriadaEm,
			},
		)	
	})

	t.Run("Fail when database connection is closed", func(t *testing.T) {
		suite.db.Close()


		err := suite.publiRepo.LikePublication(2)

		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sql: database is closed")
	})
}

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
	fmt.Println("Cleaning database...")
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
