package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

// struct que vai receber o nosso banco de dados. a lógica é que a conexão é aberta no controller e repassada para o repository realizar as manipulações no db.
type Usuarios struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

// método create do repositório de usuarios. ele recebe como parametro um modelo de usuarios e retorna um uint64 com o id do usuario inserido e um erroo.
func (u Usuarios) Create(usuario models.User) (uint64, error) {
	stmt, erro := u.db.Prepare("insert into usuarios (nome, nick, email, senha) values(?, ?, ?, ?)")

	if erro != nil {
		return 0, erro
	}

	defer stmt.Close()

	res, erro := stmt.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	lastID, erro := res.LastInsertId()

	if erro != nil {
		return 0, erro
	}

	return uint64(lastID), nil
}

// traz todos os usuarios que atendem a um filtro de nome ou nick.
func (u Usuarios) GetAll(nameOrNick string) ([]models.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick)

	if nameOrNick == "%%" {
		nameOrNick = ""
	}

	var users []models.User

	if len(nameOrNick) == 0 {
		lines, erro := u.db.Query("select id, nome, nick, email from usuarios")

		if erro != nil {
			return nil, erro
		}

		defer lines.Close()

		for lines.Next() {
			var user models.User

			if erro = lines.Scan(&user.Id, &user.Nome, &user.Nick, &user.Email); erro != nil {
				return nil, erro
			}

			users = append(users, user)
		}

		return users, nil
	} else {
		// vai retornar todos os usuarios que parcialmente forem iguais ao parametro passado para a rota, seja no campo nick, seja no campo name.
		lines, erro := u.db.Query("select id, nome, nick, email from usuarios where nome like ? or nick like ?", nameOrNick, nameOrNick)

		if erro != nil {
			return nil, erro
		}

		defer lines.Close()

		// mesma logica de cima. vai inicializar uma variavel de slice de users vazia, e depois iterar sobre as linhas retornadas pela query mysql e armazenar na variavel acima cada um dos users que é retornado. para isso, criamos uma var aux no escopo do for que a cada iteração vai receber um novo dado de usuario e depois ser repassado para o slice com o append.
		for lines.Next() {
			var user models.User

			if erro = lines.Scan(&user.Id, &user.Nome, &user.Nick, &user.Email); erro != nil {
				return nil, erro
			}

			users = append(users, user)
		}
	}

	return users, nil
}

func (u Usuarios) GetById(id uint64) (models.User, error) {
	lines, erro := u.db.Query("select id, nome, nick, email from usuarios where id = ?", id)

	// nesse caso, não conseguimos mandar um nil, porque não estamos retornando um slice como acima. por retornamos um models.User, precisamos retornar o seu valor 0, que é o model de user vazio.
	if erro != nil {
		return models.User{}, erro
	}

	defer lines.Close()

	var user models.User

	// a logica é similiar à de cima. só muda que agora não precisamos iterar sobre um slice e vamos simplesmente jogar os valores retornados pela query dentro da variavel acima, se tiver alguma linha a ser lida pelo Next().
	if lines.Next() {
		if erro = lines.Scan(&user.Id, &user.Nome, &user.Nick, &user.Email); erro != nil {
			return models.User{}, erro
		}
	} else {
		return models.User{}, fmt.Errorf("no user found with ID: %d", id)
	}

	return user, nil
}

func (u Usuarios) Update(id uint64, user models.User) error {
	_, err := u.GetById(id)
	if err != nil {
		return err
	}

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

func (u Usuarios) Delete(id uint64) error {
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

func (u Usuarios) GetByEmail(email string) (models.User, error) {
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
	} else {
		return models.User{}, fmt.Errorf("user with email %s not found", email)
	}
	return user, nil
}

func (u Usuarios) Follow(followedId uint64, followerId uint64) error {
	// o ignore vai impedir que ocorra erro ao tentar seguir um usuário que já é seguido. ou seja, se vc chamar a mesma rota duas vezes, apontando para um usuario que vc ja segue, o ignore não vai deixar ela ser executada. isso economiza processamento, pq vc n precisa fazer toda uma busca no db pra ver se determinado user ja sergue outro, e ai decidir se deixa ou n sele seguir a pessoa.
	statement, erro := u.db.Prepare("insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	// QUEM SEGUE É O USER DO TOKEN, E QUEM É SEGUIDO É O USER DA ROTA. ususario_id === seguido.
	if _, erro = statement.Exec(followerId, followedId); erro != nil {
		return erro
	}

	return nil
}

func (u Usuarios) Unfollow(followedId uint64, followerId uint64) error {
	// vai deletar a linha do DB que tem o usuario_id e o seguidor_id como seus dados.
	statement, erro := u.db.Prepare("delete from seguidores where usuario_id = ? and seguidor_id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(followerId, followedId); erro != nil {
		return erro
	}

	return nil
}

func (u Usuarios) GetFollowers(userId uint64) ([]models.User, error) {
	lines, erro := u.db.Query(`select u.id, u.nome, u.nick, u.email from usuarios u 
	inner join seguidores s on u.id = s.seguidor_id where s.usuario_id = ?`, userId)

	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var users []models.User
	for lines.Next() {
		var user models.User

		if erro = lines.Scan(
			&user.Id,
			&user.Nome,
			&user.Nick,
			&user.Email,
		); erro != nil {
			return nil, erro
		}
		users = append(users, user)
	}

	return users, nil
}

func (u Usuarios) GetFollowing(userId uint64) ([]models.User, error) {
	lines, erro := u.db.Query(`select u.id, u.nome, u.nick, u.email from usuarios u 
	inner join seguidores s on u.id = s.usuario_id where s.seguidor_id = ?`, userId)
	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var users []models.User
	for lines.Next() {
		var user models.User
		if erro = lines.Scan(
			&user.Id,
			&user.Nome,
			&user.Nick,
			&user.Email,
		); erro != nil {
			return nil, erro
		}
		users = append(users, user)
	}

	return users, nil
}

func (u Usuarios) GetPasswordFromDb(userId uint64) (string, error) {
	lines, erro := u.db.Query("select senha from usuarios where id = ?", userId)
	if erro != nil {
		return "", erro
	}

	defer lines.Close()

	var user models.User
	// para dar o scan, eu preciso instanciar um novo usuario.
	if lines.Next() {
		if erro = lines.Scan(&user.Senha); erro != nil {
			return "", erro
		}
	}

	return user.Senha, nil
}

func (u Usuarios) ChangePassword(userId uint64, newPassword string) error {
	statement, erro := u.db.Prepare("update usuarios set senha = ? where id = ?")
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(newPassword, userId); erro != nil {
		return erro
	}

	return nil
}
