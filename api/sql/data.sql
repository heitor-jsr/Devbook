insert into usuarios (nome, nick, email, senha)
values
("John Doe", "johndoe", "johndoe@example.com", "$2a$10$zMai6Qa5Gyoz3u2HiNWQ7uUOG2IvxttyVITFoBCuj0Lem794QluHG"),
("John Doe The Second", "johndoe2", "johndoe2@example.com", "$2a$10$zMai6Qa5Gyoz3u2HiNWQ7uUOG2IvxttyVITFoBCuj0Lem794QluHG"),
("John Doe The Third", "johndoe3", "johndoe3@example.com", "$2a$10$zMai6Qa5Gyoz3u2HiNWQ7uUOG2IvxttyVITFoBCuj0Lem794QluHG");

insert into seguidores(usuario_id, seguidor_id)
values
(1, 2),
(3, 1),
(1, 3);

insert into publicacoes(titulo, conteudo, autor_id)
values
("Publicação do Usuário 1", "Essa é a publicação do usuário 1! Oba!", 1),
("Publicação do Usuário 2", "Essa é a publicação do usuário 2! Oba!", 2),
("Publicação do Usuário 3", "Essa é a publicação do usuário 3! Oba!", 3);