package models

import (
	"errors"
	"strings"
	"time"
)

type Publication struct {
	Id uint64 `json:"id,omitempty"`
	Title string `json:"titulo,omitempty"`
	Content string `json:"conteudo,omitempty"`
	AuthorId uint64 `json:"autorId,omitempty"`
	AuthorNick string `json:"autorNick,omitempty"`
	Likes uint64 `json:"likes"`
	CriadaEm time.Time `json:"criadaEm,imitempty"`
}

func (publication *Publication) Prepare() error {
	if erro := publication.validate(); erro != nil {
		return erro
	}

	publication.format()
	return nil
}

func (publication *Publication) validate() error {
	if publication.Title == "" {
		return errors.New("O tiúlo da publicação é obrigatório")
	}
	if publication.Content == "" {
		return errors.New("O conteudo da publicação é obrigatório")
	}
	return nil
}

func (publication *Publication) format() {
	publication.Title = strings.TrimSpace(publication.Title)
	publication.Content = strings.TrimSpace(publication.Content)
}