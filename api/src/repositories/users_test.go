package repositories

import (
	"api/src/models"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setupTestDB() (*sql.DB, error) {

	dsn := "devbook_test:test@tcp(localhost:3306)/devbook_test?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestMain(m *testing.M) {
	db, mock := NewMock()

	sql.OpenDB("mysql", db)

	defer db.Close()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCreateUser(t *testing.T) {
	db, erro := setupTestDB()
	if erro != nil {
		t.Fatal(erro)
	}

	repository := NewUsersRepository(db)

	user := models.User{
		Nome:  "TestUser11",
		Nick:  "testuser11",
		Email: "testuser11@example.com",
		Senha: "password",
	}

	id, err := repository.Create(user)
	if err != nil {
		t.Fatal(err)
	}

	if id <= 0 {
		t.Errorf("Esperava um ID maior que zero, mas obteve %d", id)
	}
}

func TestGetAllUsers(t *testing.T) {
	db, _ := setupTestDB()
	defer db.Close()

	repository := NewUsersRepository(db)

	_, err := repository.GetAll("TestUser")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserByID(t *testing.T) {
	db, _ := setupTestDB()
	defer db.Close()

	repository := NewUsersRepository(db)

	user, err := repository.GetById(1)
	if err != nil {
		t.Fatal(err)
	}

	if user.Id != 1 {
		t.Errorf("Esperava um ID de usuário 1, mas obteve %d", user.Id)
	}
}
