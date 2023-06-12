package postgres

const (
	CREATE_TABLE_STATUS = `CREATE TABLE status (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		code_name TEXT NOT NULL
	);`

	CREATE_TABLE_USERS = `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash CHAR(60) NOT NULL
	);`
)