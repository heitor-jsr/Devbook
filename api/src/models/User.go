package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
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
func (user *User) Prepare(step string) error {
	if erro := user.validate(step); erro != nil {
		return erro
	}
	if erro := user.format(step); erro != nil {
		return erro
	}
	return nil
}

//verifica se todos os campos do struct de user tão preenchidos.
func (user *User) validate(step string) error {
	if user.Nome == "" {
		return errors.New("o campo nome é obrigatorio")
	}	
	if step == "create" && user.Senha == "" {
		return errors.New("o campo senha é obrigatorio")
	}
	if user.Nick == "" {
		return errors.New("o campo nick é obrigatorio")
	}
	if user.Email == "" {
		return errors.New("o campo email é obrigatorio")
	}
	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("o email inserido é invalido")
	}
	return nil
}

// formata os campos do struct para que não exista espaço em branco nas extremidades dos dados inseridos.
func (user *User) format(step string) error {
	user.Nome = strings.TrimSpace(user.Nome)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if step == "create" {
		hashedPassword, erro := security.Hash(user.Senha)
		if erro != nil {
			return erro
		}

		user.Senha = string(hashedPassword)
	}
	return nil
}