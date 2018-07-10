package db

import (
	"database/sql"

	"github.com/cohhei/go-to-the-handson/04/schema"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func ConnectPostgres() (*Postgres, error) {
	connStr := "postgres://postgres:postgres@postgres/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (p *Postgres) Close() {
	p.DB.Close()
}

func (p *Postgres) Insert(todo *schema.Todo) (int, error) {
	query := `
		INSERT INTO todo (id, title, note, due_date)
		VALUES (nextval('todo_id'), $1, $2, $3)
		RETURNING id;
	`

	rows, err := p.DB.Query(query, todo.Title, todo.Note, todo.DueDate)
	if err != nil {
		return -1, err
	}

	var id int
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (p *Postgres) Delete(id int) error {
	query := `
		DELETE FROM todo
		WHERE id = $1;
	`

	if _, err := p.DB.Exec(query, id); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetAll() ([]schema.Todo, error) {
	query := `
		SELECT *
		FROM todo
		ORDER BY id;
	`

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var todoList []schema.Todo
	for rows.Next() {
		var t schema.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Note, &t.DueDate); err != nil {
			return nil, err
		}
		todoList = append(todoList, t)
	}

	return todoList, nil
}
