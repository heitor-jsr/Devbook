package models

import (
	"time"
	"errors"
	"strings"
)

type User struct {
	Id uint64 `json:"id,omitempty"`
	Nome string `json:"nome,omitempty"`
	Nick string `json:"nick,omitempty"`
	Email string `json:"email,omitempty"`
	Senha string `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoEm,omitempty"`
}

// chama os métodos abaixo, verificando se eles são válidos e, se não, retorna erro. se forem válidos, os dados são formatados e a execução d o código segue.
func (user *User) Prepare() error {
	if erro := user.validate(); erro != nil {
		return erro
	}
	user.format()
	return nil
}

//verifica se todos os campos do struct de user tão preenchidos.
func (user *User) validate() error {
	if user.Nome == "" {
		return errors.New("o campo nome é obrigatorio")
	}	
	if user.Senha == "" {
		return errors.New("o campo senha é obrigatorio")
	}
	if user.Nick == "" {
		return errors.New("o campo nick é obrigatorio")
	}
	if user.Email == "" {
		return errors.New("o campo email é obrigatorio")
	}
	return nil
}

// formata os campos do struct para que não exista espaço em branco nas extremidades dos dados inseridos.
func (user *User) format() {
	user.Nome = strings.TrimSpace(user.Nome)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)
}