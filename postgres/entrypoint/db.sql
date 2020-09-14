\connect postgres

CREATE DATABASE avitojob;

\connect avitojob

CREATE TABLE Users (
	user_id serial NOT NULL,
	balance double precision NOT NULL,
	CONSTRAINT Users_pk PRIMARY KEY (user_id)
) WITH (
  OIDS=FALSE
);



CREATE TABLE Transactions (
	trans_id serial NOT NULL,
	user_id integer NOT NULL,
	init_balance double precision NOT NULL,
	change double precision NOT NULL,
	time VARCHAR(255) NOT NULL,
	source VARCHAR(255) NOT NULL,
	comment VARCHAR(255) NOT NULL,
	CONSTRAINT Transactions_pk PRIMARY KEY (trans_id)
) WITH (
  OIDS=FALSE
);


ALTER TABLE Transactions ADD CONSTRAINT Transactions_fk0 FOREIGN KEY (user_id) REFERENCES Users(user_id);

