package repositories

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func StartMySQLContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string, func() (string, error)) {
	createTablesScriptPath := filepath.Join("..", "..", "sql", "create_tables.sql")

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8"),
		mysql.WithDatabase("devbook"),
		mysql.WithUsername("golang"),
		mysql.WithPassword("golang"),
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

	dsn := fmt.Sprintf("golang:golang@tcp(%s:%s)/devbook", ip, port.Port())

	teardown := func() (string, error) {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			return "", err
		}
		return dsn, nil
	}

	return mysqlContainer, dsn, teardown
}
