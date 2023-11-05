package repositories

import (
	"api/src/models"
	"database/sql"
)

// struct que vai receber o nosso banco de dados. a lógica é que a conexão é aberta no controller e repassada para o repository realizar as manipulações no db.
type publications struct {
	db *sql.DB
}

func NewPublicationRepository(db *sql.DB) *publications {
	return &publications{db}
}

func (publications *publications) Create(publication models.Publication) (uint64, error) { 
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

func (publications *publications) GetPublicationById(publicationId uint64) (models.Publication, error) { 
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

	var publication models.Publication
	if lines.Next() {
		if erro = lines.Scan(
			&publication.Id,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&publication.CriadaEm,
			&publication.AuthorNick,
		); erro != nil {
			return models.Publication{}, erro
		}
	}

	return publication, nil
}