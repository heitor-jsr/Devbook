package repositories

import (
	"api/src/models"
	"database/sql"
	"time"
)

// struct que vai receber o nosso banco de dados. a lógica é que a conexão é aberta no controller e repassada para o repository realizar as manipulações no db.
type Publications struct {
	db *sql.DB
}

func NewPublicationRepository(db *sql.DB) *Publications {
	return &Publications{db}
}

func (publications *Publications) Create(publication models.Publication) (uint64, error) { 
	statement, erro := publications.db.Prepare(
		"insert into publicacoes (titulo, conteudo, autor_id) values (?, ?, ?)")
	if erro != nil {
		return 0, erro
	}

	defer statement.Close()

	result, erro := statement.Exec(publication.Title, publication.Content, publication.AuthorId)
	if erro != nil {
		return 0, erro
	}

	publicationId, erro := result.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(publicationId), nil
}

func (publications *Publications) GetPublicationById(publicationId uint64) (models.Publication, error) { 
	lines, erro := publications.db.Query(
		`select p.*, u.nick from 
		publicacoes p inner join usuarios u
		on u.id = p.autor_id where p.id = ?`,
		publicationId,
	)

	if erro != nil {
		return models.Publication{}, erro
	}

	defer lines.Close()

	var criadaEmStr string
	var publication models.Publication
	
	if lines.Next() {
		if erro = lines.Scan(
			&publication.Id,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&criadaEmStr,
			&publication.AuthorNick,
		); erro != nil {
			return models.Publication{}, erro
		}
		publication.CriadaEm, erro = time.Parse("2006-01-02 15:04:05", criadaEmStr)
		if erro != nil {
			return models.Publication{}, erro
		}
	}

	return publication, nil
}

// vai retornar todas as publicações dos usuários que ele segue e todas as publicações pŕoprias dele.
func (publications *Publications) GetPublications(usuarioID uint64) ([]models.Publication, error) { 
	lines, erro := publications.db.Query(`
	select distinct p.*, u.nick from publicacoes p 
	inner join usuarios u on u.id = p.autor_id 
	inner join seguidores s on p.autor_id = s.usuario_id 
	where u.id = ? or s.seguidor_id = ?
	order by 1 desc`,
		usuarioID, usuarioID,
	)

	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var newPublications []models.Publication
	var criadaEmStr string
	for lines.Next() {
		var publication models.Publication
		if erro = lines.Scan(
			&publication.Id,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&criadaEmStr,
			&publication.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		publication.CriadaEm, erro = time.Parse("2006-01-02T15:04:05Z", criadaEmStr)
    if erro != nil {
        return nil, erro
    }
		newPublications = append(newPublications, publication)
	}

	return newPublications, nil
}

func (publications *Publications) Update(publicationId uint64, pupublication models.Publication) (error) { 
	statement, erro := publications.db.Prepare(
		"update publicacoes set titulo = ?, conteudo = ? where id = ?")
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(pupublication.Title, pupublication.Content, publicationId); erro != nil {
		return erro
	}

	return nil
}

func (publications *Publications) Delete(publicationId uint64) (error) { 
	statement, erro := publications.db.Prepare(
		"delete from publicacoes where id = ?",
	)
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(publicationId); erro != nil {
		return erro
	}

	return nil
}

func (publications *Publications) GetPublicationFromUser(userId uint64) ([]models.Publication, error) {
	lines, erro := publications.db.Query(`
	select p.*, u.nick from publicacoes p
	inner join usuarios u on u.id = p.autor_id
	where p.autor_id = ?
	order by 1 desc`,
		userId,
	)
	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var newPublications []models.Publication
	var criadaEmStr string

	for lines.Next() {
		var publication models.Publication
		if erro = lines.Scan(
			&publication.Id,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&criadaEmStr,
			&publication.AuthorNick,
		); erro != nil {
			return nil, erro
		}
		publication.CriadaEm, erro = time.Parse("2006-01-02T15:04:05Z", criadaEmStr)
		if erro != nil {
			return nil, erro
		}
		newPublications = append(newPublications, publication)
	}

	return newPublications, nil
}

func (publications *Publications) LikePublication(publicationId uint64) (error) {
	statement, erro := publications.db.Prepare(
		"update publicacoes set curtidas = curtidas + 1 where id = ?",
	)
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(publicationId); erro != nil {
		return erro
	}

	return nil
}

func (publications *Publications) DeslikePublication(publicationId uint64) (error) {
	statement, erro := publications.db.Prepare(
		"update publicacoes set curtidas = case when curtidas > 0 then curtidas - 1 else curtidas end where id = ?",
	)
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(publicationId); erro != nil {
		return erro
	}

	return nil
}