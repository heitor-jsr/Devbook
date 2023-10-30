package database

import (
	"api/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// o responsavel por chamar essa função e abrir a conexão com o banco de dados é o controller. depois de aberta, o controller repassa a conexão para os repositorios iterarem sobre as tabelas do banco de dados.
func Connect() (*sql.DB, error) {
	db, erro := sql.Open("mysql", config.ConnStr)
	if erro != nil {
		return nil, erro
	}

	// se teve algum problema ao abrir a conexão com o banco, precisamos garantir que ela seja fechada.
	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}