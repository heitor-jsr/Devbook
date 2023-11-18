insert into usuarios (nome, nick, email, senha)
values
("John Doe", "johndoe", "johndoe@example.com", "$2a$10$zMai6Qa5Gyoz3u2HiNWQ7uUOG2IvxttyVITFoBCuj0Lem794QluHG"),
("John Doe The Second", "johndoe2", "johndoe2@example.com", "$2a$10$zMai6Qa5Gyoz3u2HiNWQ7uUOG2IvxttyVITFoBCuj0Lem794QluHG");

insert into publicacoes(titulo, conteudo, autor_id)
values
("John Doe Publication", "John Doe Publication Content", 1),
("John Doe The Second", "John Doe The Second Content", 2);


insert into seguidores(usuario_id, seguidor_id)
values
(1, 2),
(2, 1)
