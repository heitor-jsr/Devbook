package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

// struct que vai receber o nosso banco de dados. a lógica é que a conexão é aberta no controller e repassada para o repository realizar as manipulações no db.
type usuarios struct {
	db *sql.DB
}

// função que recebe um banco aberto pelo controller como parametro, sendo o controller que chama essa func trambém. a função, por sua vez, pega o banco e joga dentro do struct de usuários que criamos. ou seja, instanciamos o struct com o banco aberto pelo controller.
func NewUsersRepository(db *sql.DB) *usuarios {
	return &usuarios{db}
}

// método create do repositório de usuarios. ele recebe como parametro um modelo de usuarios e retorna um uint64 com o id do usuario inserido e um erro.
func (u usuarios) Create(usuario models.User) (uint64, error) {
	stmt, err := u.db.Prepare("insert into usuarios (nome, nick, email, senha) values(?, ?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}