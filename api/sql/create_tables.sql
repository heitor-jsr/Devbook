CREATE DATABASE IF NOT EXISTS devbook;

DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS seguidores;
DROP TABLE IF EXISTS publicacoes;

CREATE TABLE usuarios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome VARCHAR(50) NOT NULL CHECK (nome <> ''),
  nick VARCHAR(50) NOT NULL UNIQUE CHECK (nick <> ''),
  email VARCHAR(50) NOT NULL UNIQUE CHECK (email <> ''),
  senha VARCHAR(250) NOT NULL CHECK (senha <> ''),
  criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
) ENGINE=INNODB;

CREATE TABLE seguidores (
  usuario_id INT NOT NULL, FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
  seguidor_id INT NOT NULL, FOREIGN KEY (seguidor_id) REFERENCES usuarios(id) ON DELETE CASCADE,
  PRIMARY KEY (usuario_id, seguidor_id)
) ENGINE=INNODB;

CREATE TABLE publicacoes(
    id int auto_increment primary key,
    titulo varchar(50) NOT NULL CHECK (titulo <> ''),
    conteudo varchar(300) NOT NULL CHECK (conteudo <> ''),

    autor_id int NOT NULL CHECK (autor_id <> ''),
    FOREIGN KEY (autor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    curtidas int default 0,
    criadaEm timestamp default current_timestamp
) ENGINE=INNODB;
