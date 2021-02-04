CREATE DATABASE IF NOT EXISTS ginopexample;

USE ginopexample;

CREATE TABLE IF NOT EXISTS book (
    name        VARCHAR(100),
    author       VARCHAR(100),
    genre VARCHAR(100),
    PRIMARY KEY (name)
);

INSERT INTO book VALUES ('A Game of Thrones', 'George RR Martin', 'fantasy');
INSERT INTO book VALUES ('A Play of Giants', 'Wole Soyinka', 'drama');
INSERT INTO book VALUES ('The Famished Road', 'Ben Okri', 'fiction');
INSERT INTO book VALUES ('Animal Farm', 'Geroge Orwell', 'fiction');
INSERT INTO book VALUES ('The Importance of Being Earnest', 'Oscar Wilder', 'drama');