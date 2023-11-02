package auth

import (
	"api/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// obs: o ideal é gerar o secret e armazena-lo dentro do .env.

func GenerateToken(userID uint64) (string, error) {
	// map que contém as permissões do usuario autenticado.
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userId"] = userID
	// essas são as permissoes dentro do nosso token. a primeira coisa a se fazer é gerar ele. o proximo passo é assinar digitalmente esse token. para assinar, usamos a chave secret. fazemos isso com a função abaixo. o primeiro metodo é como iremos assinar o token, e o segundo são as permissões do token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	// criado o token com a func acima, basta assina-lo com o secret. o secret é a chave responsavel por fazer a assinatura do token e garantir a autenticidade dele. é recomendado q ela seja gerada de forma segura.
	return token.SignedString([]byte(config.SecretKey))

	// o token vai conter a identificação do usuario a que ele se refere, as autorizações dele, quando ele expira, etc.
}