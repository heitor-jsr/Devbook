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

func (u usuarios) GetById(id uint64) (models.User, error) {
	lines, erro := u.db.Query("select id, nome, nick, email, criadoEm from usuarios where id = ?", id)

	// nesse caso, não conseguimos mandar um nil, porque não estamos retornando um slice como acima. por retornamos um models.User, precisamos retornar o seu valor 0, que é o model de user vazio.
	if erro != nil {
		return models.User{}, erro
	}

	defer lines.Close()

	var user models.User

	// a logica é similiar à de cima. só muda que agora não precisamos iterar sobre um slice e vamos simplesmente jogar os valores retornados pela query dentro da variavel acima, se tiver alguma linha a ser lida pelo Next().
	if lines.Next() {
		if erro = lines.Scan(&user.Id, &user.Nome, &user.Nick, &user.Email, &user.CriadoEm); erro != nil {
			return models.User{}, erro
		}
	}

	return user, nil
}

func (u usuarios) Update(id uint64 , user models.User) error {
	statemente, erro := u.db.Prepare("update usuarios set nome = ?, nick = ?, email = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statemente.Close()

	// executa o statemente com os valores a serem alterados na query. lembrando que o prepare é usado para lidar com o sql injection e só é utilizado nos casos em que há manipulação dos dados nas tabelas sql.
	if _, erro = statemente.Exec(user.Nome, user.Nick, user.Email, id); erro != nil {
		return erro
	}

	return nil
}

func (u usuarios) Delete(id uint64) error {
	statemente, erro := u.db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		return erro
	}
	defer statemente.Close()

	if _, erro = statemente.Exec(id); erro != nil {
		return erro
	}

	return nil
}

func (u usuarios) GetByEmail(email string) (models.User, error) {
	lines, erro := u.db.Query("select id, senha from usuarios where email = ?", email)
	if erro != nil {
		return models.User{}, erro
	}
	defer lines.Close()

	var user models.User
	if lines.Next() {
		// scanea as linhas retornadas pelo query e armazena nas variaveis auxiliares criadas acima, no endereço de memoria delas.
		if erro = lines.Scan(&user.Id, &user.Senha); erro != nil {
			return models.User{}, erro
		}
	}
	return user, nil
}

func (u usuarios) Follow(userId uint64, followerId uint64) error {
	// o ignore vai impedir que ocorra erro ao tentar seguir um usuário que já é seguido. ou seja, se vc chamar a mesma rota duas vezes, apontando para um usuario que vc ja segue, o ignore não vai deixar ela ser executada. isso economiza processamento, pq vc n precisa fazer toda uma busca no db pra ver se determinado user ja sergue outro, e ai decidir se deixa ou n sele seguir a pessoa.
	statement, erro := u.db.Prepare("insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(userId, followerId); erro != nil {
		return erro
	}

	return nil
}

func (u usuarios) Unfollow(userId uint64, followerId uint64) error {
	// vai deletar a linha do DB que tem o usuario_id e o seguidor_id como seus dados.
	statement, erro := u.db.Prepare("delete from seguidores where usuario_id = ? and seguidor_id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(userId, followerId); erro != nil {
		return erro
	}

	return nil
}