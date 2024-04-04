/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE account (
	id serial PRIMARY KEY,
	full_name VARCHAR ( 60 ) NOT NULL,
    phone VARCHAR ( 13 ) UNIQUE NOT NULL,
    password VARCHAR ( 64 ) NOT NULL
);

CREATE TABLE login (
    id serial PRIMARY KEY,
    account_id serial UNIQUE,
    success_login INT DEFAULT 0,
    token VARCHAR NOT NULL,
    CONSTRAINT fk_account_auth FOREIGN KEY(account_id) REFERENCES account(id)
);
