package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func SetupMySQLContainer() (*sql.DB, error) {
	ctx := context.Background()

	sqlScriptPath := filepath.Join("..", "..", "sql", "sql.sql")

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8"),
		mysql.WithDatabase("devbook"),
		mysql.WithUsername("golang"),
		mysql.WithPassword("golang"),
		mysql.WithScripts(sqlScriptPath),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
				panic(err)
		}
	}()

	containerPort, err := mysqlContainer.MappedPort(ctx, "3306/tcp")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	dsn := fmt.Sprintf("golang:golang@tcp(127.0.0.1:%s)/devbook", containerPort.Port())

	db, err := sql.Open("mysql", dsn)
	if err != nil {
			log.Fatal(err)
	}
	return db, nil
}