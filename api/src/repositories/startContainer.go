package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func StartMySQLContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string, func() (string, error)) {
	createTablesScriptPath := filepath.Join("..", "..", "sql", "create_tables.sql")

	if err := godotenv.Load("/home/dornzak/devbook/api/.env"); err != nil {
		log.Fatal(err)
	}

	mysqlUser :=	os.Getenv("DB_USER")
	mysqlPassword := os.Getenv("DB_PASSWORD")
	mysqlDbName := os.Getenv("DB_NAME")

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8"),
		mysql.WithDatabase(mysqlDbName),
		mysql.WithUsername(mysqlUser),
		mysql.WithPassword(mysqlPassword),
		mysql.WithScripts(createTablesScriptPath),
	)
	if err != nil {
		t.Fatal(err)
	}

	ip, err := mysqlContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306/tcp")
	if err != nil {
		t.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, ip, port.Port(), mysqlDbName)
	
	teardown := func() (string, error) {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			return "", err
		}
		return dsn, nil
	}

	return mysqlContainer, dsn, teardown
}
