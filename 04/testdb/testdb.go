package testdb

import (
	"database/sql"

	"github.com/corhhey/go-to-the-handson/04/db"
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

func Setup() *db.Postgres {
	sqlDB, err := connectPostgresForTests()
	if err != nil {
		panic(err)
	}

	if _, err = sqlDB.Exec(createTable); err != nil {
		panic(err)
	}

	return &db.Postgres{sqlDB}
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
