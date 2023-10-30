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

// método create do repositório de usuarios. ele recebe como parametro um modelo de usuarios e retorna um uint64 com o id do usuario inserido e um erroo.
func (u usuarios) Create(usuario models.User) (uint64, error) {
	stmt, erro := u.db.Prepare("insert into usuarios (nome, nick, email, senha) values(?, ?, ?, ?)")

	if erro != nil {
		return 0, erro
	}

	defer stmt.Close()

	res, erro := stmt.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		fmt.Println(erro)
		return 0, erro
	}

	lastID, erro := res.LastInsertId()

	if erro != nil {
		return 0, erro
	}

	return uint64(lastID), nil
}