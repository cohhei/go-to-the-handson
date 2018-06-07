package db

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/corhhey/go-to-the-handson/04/schema"
	_ "github.com/lib/pq"
)

func TestPostgres_Insert(t *testing.T) {
	postgres := setup(t)
	defer postgres.Close()

	todo := &schema.Todo{
		Title:   "title1",
		Note:    "note1",
		DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
	}

	got, err := postgres.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	want := 1

	if got != want {
		t.Fatal(err)
	}
}

func TestPostgres_GetAll(t *testing.T) {
	postgres := setup(t)
	defer postgres.Close()

	todo := &schema.Todo{
		Title:   "title1",
		Note:    "note1",
		DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
	}

	_, err := postgres.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	got, err := postgres.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	want := []schema.Todo{
		{
			ID:      1,
			Title:   "title1",
			Note:    "note1",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
	}

	if equal(got, want) {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}

func TestPostgres_Delete(t *testing.T) {
	postgres := setup(t)
	defer postgres.Close()

	todo := &schema.Todo{
		Title:   "title1",
		Note:    "note1",
		DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
	}

	id, err := postgres.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	err = postgres.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	got, err := postgres.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) > 0 {
		t.Fatal("The record is not deleted.")
	}
}

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

func setup(t *testing.T) *Postgres {
	db, err := connectPostgresForTests()
	if err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec(createTable); err != nil {
		t.Fatal(err)
	}

	return &Postgres{db}
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

func equal(got interface{}, want interface{}) bool {
	return reflect.DeepEqual(got, want)
}
