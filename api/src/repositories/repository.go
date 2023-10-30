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

// traz todos os usuarios que atendem a um filtro de nome ou nick.
func (u usuarios) GetAll(nameOrNick string) ([]models.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick)

	// vai retornar todos os usuarios que parcialmente forem iguais ao parametro passado para a rota, seja no campo nick, seja no campo name.
	lines, erro := u.db.Query("select id, nome, nick, email, criadoEm from usuarios where nome like ? or nick like ?", nameOrNick, nameOrNick)
	
	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var users []models.User
	// mesma logica de cima. vai inicializar uma variavel de slice de users vazia, e depois iterar sobre as linhas retornadas pela query mysql e armazenar na variavel acima cada um dos users que é retornado. para isso, criamos uma var aux no escopo do for que a cada iteração vai receber um novo dado de usuario e depois ser repassado para o slice com o append.
	for lines.Next() {
		var user models.User

		if erro = lines.Scan(&user.Id, &user.Nome, &user.Nick, &user.Email, &user.CriadoEm); erro != nil {
			return nil, erro
		}

		users = append(users, user)
	}

	return users, nil
}
