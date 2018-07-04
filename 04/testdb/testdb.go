package testdb

import (
	"database/sql"
)

const createTable = `
DROP TABLE IF EXISTS todo;
Alter SEQUENCE todo_id RESTART WITH 1;
CREATE TABLE todo (
  ID serial PRIMARY KEY,
  TITLE TEXT NOT NULL,
  NOTE TEXT,
  DUE_DATE TIMESTAMP WITH TIME ZONE
);
`

type TestDB struct {
	db *sql.DB
}

func Setup() *sql.DB {
	db, err := connectPostgresForTests()
	if err != nil {
		panic(err)
	}

	if _, err = db.Exec(createTable); err != nil {
		panic(err)
	}

	return db
}

func connectPostgresForTests() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
