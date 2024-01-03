#!/bin/sh

# Aguarde o MySQL estar pronto
/wait-for-it.sh mysql:3306 -timeout=60 -- echo "MySQL está pronto, iniciando a aplicação."

# Inicie sua aplicação
exec /usr/local/bin/devbook/api
