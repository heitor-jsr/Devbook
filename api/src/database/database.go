package database

import (
	"api/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// o responsavel por chamar essa função e abrir a conexão com o banco de dados é o controller. depois de aberta, o controller repassa a conexão para os repositorios iterarem sobre as tabelas do banco de dados.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.ConnStr)
	if err != nil {
		return nil, err
	}

	// se teve algum problema ao abrir a conexão com o banco, precisamos garantir que ela seja fechada.
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}