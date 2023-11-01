package security

import "golang.org/x/crypto/bcrypt"

func Hash(password string) ([]byte, error) {
	// vai gerar um hash a partir de uma senha. o primeiro parametro é a senha que é bassada como parametro da func em slice de byte, e o segundo é o custo da operação (numero de vezes que, conforme aumenta, aumenta a complexidade do hash).
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	// compara se a senha que é passada como parametro, e a senha que foi gerada pelo hash, sao iguais. se forem iguais, o retorno é nulo. caso contrário, o retorno é um erro.
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}