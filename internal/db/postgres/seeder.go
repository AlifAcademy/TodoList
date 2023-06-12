package postgres

const (
	CreateTableStatus = `CREATE TABLE status (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		code_name TEXT NOT NULL
	  );`
	
	  CreateTableUSERS = `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash CHAR(60) NOT NULL
	  );`
	
	  CreateTableTasks = `CREATE TABLE tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		tags TEXT[],
		status_id INT NOT NULL REFERENCES status(id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
	  );`
	
	  CreateTableComments = `CREATE TABLE comments (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
	  );`
)