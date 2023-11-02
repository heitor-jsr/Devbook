package auth

import (
	"api/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	// para verificar de fato o token, precisamos dar um parse nele para retirar as autorizações d o usuarios e verificar se estamos mesmo lidando com um token válido. o primeiro parametro da func é a string do token, e a segunda é uma função que vai retornar a chave de verificação que devemos usar pra dar o parse no token. precisamos retornar essa chave pq o jwt orienta que, antes de dar o parse, a gente precisa verificar se o método de assinatura do token que vamos dar o parse é o que estamos esperando, porque não podemos assinar o token com um método e dar o parse nele usando um outro método de assinatura.
	token, erro := jwt.Parse(tokenString, returnVerifyKey)
	if erro != nil {
		return erro
	}

	// vai verificar se os claims do meu token existem e se o token ainda não expirou.
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("token invalido")
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	// vc precisa dar um split nesse token pq quando vc pega ele da req, ele vem assim: Bearer <token>, e só precisamos do token nesse caso.
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

// Vai retornar se o método de assinatura que estamos usando para o caso é de uma familia específica.
func returnVerifyKey(token *jwt.Token) (interface{}, error) {
	// o signingMethodHMC representa uma família de métodos de assinatura.
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// se der erro, informamos que o método de assinatura esperado é incorreto, e informamos o método que foi tentado.
		return nil, fmt.Errorf("metodo de assinatura inesperado! %v", token.Header["alg"])
	}
	// o token que estamos passando e o secret do .env.
	return config.SecretKey, nil
}