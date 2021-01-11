CREATE DATABASE IF NOT EXISTS ginopexample;

USE ginopexample;

CREATE TABLE IF NOT EXISTS book (
    name        VARCHAR(100),
    author       VARCHAR(10),
    genre VARCHAR(100),
    PRIMARY KEY (name)
);

INSERT INTO book VALUES ("A Game of Thrones", "George RR Martin", "Fantasy");
INSERT INTO book VALUES ("A Play of Giants", "Wole Soyinka", "Drama");
INSERT INTO book VALUES ("The Famished Road", "Ben Okri", "Fiction");
INSERT INTO book VALUES ("Animal Farm", "Geroge Orwell", "Fiction");
INSERT INTO book VALUES ("The Importance of Being Earnest", "Oscar Wilder", "Drama");