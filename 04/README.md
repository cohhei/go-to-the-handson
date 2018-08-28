# 4. Sample Application -TODO List API-

We are going to develop a tiny RESTful API for the TODO application that contains the following features:

- Return the all TODO
- Save new TODO
- Delete a TODO

Before building the API, please create new directory.

```sh
$ mkdir todo
$ cd todo
```

contents:
<!-- TOC -->

- [4. Sample Application -TODO List API-](#4-sample-application--todo-list-api-)
  - [Data schema](#data-schema)
  - [Repositories](#repositories)
  - [Sample TODO list](#sample-todo-list)
    - [Test for samples](#test-for-samples)
    - [Implementation for samples](#implementation-for-samples)
  - [Service](#service)
  - [Routing and Handlers](#routing-and-handlers)
    - [Sample Handler](#sample-handler)
    - [Routing](#routing)
  - [Functional test](#functional-test)
  - [main.go](#maingo)
  - [PostgreSQL](#postgresql)
    - [DB for tests](#db-for-tests)
    - [Test for Postgres](#test-for-postgres)
    - [Implementation for Postgres](#implementation-for-postgres)
    - [Functional test for Postgres](#functional-test-for-postgres)
    - [New routing and handlers](#new-routing-and-handlers)
    - [Updating main.go](#updating-maingo)
  - [Docker](#docker)
    - [Postgres](#postgres)
    - [TODO API](#todo-api)
    - [docker-compose.yml](#docker-composeyml)
  - [Requests](#requests)

<!-- /TOC -->

## Data schema

The API deal TODO lists with clients as a JSON files such as:

```json
[
  {
    "id": 1,
    "title": "Do dishes",
    "note": "That will be done by Gopher.",
    "due_date": "2000-01-01T00:00:00Z"
  }
]
```

First of all, let's define the schema with Go.

```sh
$ mkdir schema
$ touch schema/model.go
```

```go
// schema/model.go
package schema

import "time"

type Todo struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Note    string    `json:"note"`
	DueDate time.Time `json:"due_date"`
}
```

It's a really common pattern to define schemas as a structure. The `Todo` structure has four fields, `ID`, `Title`, `Note`, and `DueDate`. They respectively have their types and [json tags](https://golang.org/pkg/encoding/json/). The tags will be the field names when the `Todo` structure is encoded. If you don't write the tags, json field names will be the same as the structure's name.

## Repositories

Create the `db` directory and the `db/repository.go`.

```sh
$ mkdir db
$ touch db/repository.go
```

```go
// db/repository.go
package db

import (
	"context"

	"../schema"
)

const keyRepository = "Repository"

type Repository interface {
	Close()
	Insert(todo *schema.Todo) (int, error)
	Delete(id int) error
	GetAll() ([]schema.Todo, error)
}

func SetRepository(ctx context.Context, repository Repository) context.Context {
	return context.WithValue(ctx, keyRepository, repository)
}

func Close(ctx context.Context) {
	getRepository(ctx).Close()
}

func Insert(ctx context.Context, todo *schema.Todo) (int, error) {
	return getRepository(ctx).Insert(todo)
}

func Delete(ctx context.Context, id int) error {
	return getRepository(ctx).Delete(id)
}

func GetAll(ctx context.Context) ([]schema.Todo, error) {
	return getRepository(ctx).GetAll()
}

func getRepository(ctx context.Context) Repository {
	return ctx.Value(keyRepository).(Repository)
}
```

This `Repository` interface has four methods. We can access our DB through the interface in the [context](https://golang.org/pkg/context/). The interface can divide the application logic and implementations for each middleware such as PostgreSQL or MySQL. The `SetRepository` function sets structures implementing the `Repository` interface to the context.

## Sample TODO list

At first, we'll create a structure that returns static TODO list as samples instead of dynamic values in DB.

```sh
$ touch db/samples.go
```

```go
// db/samples.go
package db

import "../schema"

type Sample struct{}

func (s *Sample) Close() {}

func (s *Sample) Insert(todo *schema.Todo) (int, error) {
	return 0, nil
}

func (s *Sample) Delete(id int) error {
	return nil
}

func (s *Sample) GetAll() ([]schema.Todo, error) {
	return nil, nil
}
```

The `Sample` structure has `Close`, `Insert`, `Delete`, and `GetAll` methods, so we can use it as a `Repository` interface type. But these methods do nothing yet. Create tests for them first and implement them after that.

### Test for samples

Create `db/sample_test.go` and write the tests.

```sh
$ touch db/sample_test.go
```

```go
// db/sample_test.go
package db

import (
	"reflect"
	"testing"
	"time"

	"../schema"
)

func TestClose(t *testing.T) {
	sample := Sample{}

	sample.Close()
}

func TestInsert(t *testing.T) {
	sample := Sample{}

	todo := &schema.Todo{}

	got, err := sample.Insert(todo)
	if err != nil {
		t.Error(err)
	}

	if got != 0 {
		t.Fatal("Want: 0, Got: ", got)
	}
}

func TestDelete(t *testing.T) {
	sample := Sample{}

	err := sample.Delete(1)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAll(t *testing.T) {
	sample := Sample{}

	got, err := sample.GetAll()
	if err != nil {
		t.Error(err)
	}

	want := []schema.Todo{
		{
			ID:      1,
			Title:   "Do dishes",
			Note:    "",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      2,
			Title:   "Do homework",
			Note:    "",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      2,
			Title:   "Twitter",
			Note:    "",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want: %v, Got: %v\n", want, got)
	}
}
```

Our sample repository won't have any completed features for `Close`, `Insert`, and `Delete`. So the tests for them are meaningless and already pass. Of course, `TestGetAll` will fail yet.

```sh
$ go test ./db/*
--- FAIL: TestGetAll (0.00s)
        sample_test.go:71: Want: [{1 Do dishes  2000-01-01 00:00:00 +0000 UTC} {2 Do homework  2000-01-01 00:00:00 +0000 UTC} {2 Twitter  2000-01-01 00:00:00 +0000 UTC}], Got: []
FAIL
FAIL    command-line-arguments  0.041s
```

### Implementation for samples

```go
// db/samples.go
package db

import (
	"time"

	"../schema"
)

// ...

func (s *Sample) GetAll() ([]schema.Todo, error) {
	todoList := []schema.Todo{
		{
			ID:      1,
			Title:   "Do dishes",
			Note:    "",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      2,
			Title:   "Do homework",
			Note:    "",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      2,
			Title:   "Twitter",
			Note:    "",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	return todoList, nil
}
```

Execute the tests. They will pass.

```sh
$ go test ./db/*
ok      command-line-arguments  0.182s
```

## Service

Create the `service` directory and `service/todo.go`. This service has only simple functions that call functions with the same name in the `db` package. If you'd like to upgrade this TODO API, the service will have a lot of complex function.

```sh
$ mkdir service
$ touch service/todo.go
```

```go
// service/todo.go
package service

import (
	"context"

	"../db"
	"../schema"
)

func Close(ctx context.Context) {
	db.Close(ctx)
}

func Insert(ctx context.Context, todo *schema.Todo) (int, error) {
	return db.Insert(ctx, todo)
}

func Delete(ctx context.Context, id int) error {
	return db.Delete(ctx, id)
}

func GetAll(ctx context.Context) ([]schema.Todo, error) {
	return db.GetAll(ctx)
}
```

## Routing and Handlers

### Sample Handler

```sh
$ mkdir handler
$ touch handler/todo.go
```

```go
// handler/todo.go
package handler

import (
	"encoding/json"
	"net/http"

	"../db"
	"../service"
)

type todoHandler struct {
	samples  *db.Sample
}

func (handler *todoHandler) GetSamples(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetRepository(r.Context(), handler.samples)

	todoList, err := service.GetAll(ctx)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseOk(w, todoList)
}
```

The package have the `GetSamples` method that sets a `db.Sample` structure to the context. So the TODO service will return sample TODO list.

Create utility functions `responseOk` and `responseError`.

```go
// handler/todo.go
func responseOk(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(body)
}

func responseError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	body := map[string]string{
		"error": message,
	}
	json.NewEncoder(w).Encode(body)
}
```

### Routing

Create `handler/routes.go`. The routings will be written in the file.

```sh
$ touch handler/routes.go
```

```go
// handler/routes.go
package handler

import (
	"net/http"

	"../db"
)

func SetUpRouting() *http.ServeMux {
	todoHandler := &todoHandler{
		samples:  &db.Sample{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/samples", todoHandler.GetSamples)
	
	return mux
}
```

## Functional test

We've completed building the `/samples` endpoint the should return the sample TODO list as a Json file. Let's write the functional test to confirm it.

```sh
$ mkdir functional
$ touch functional/todo_test.go
```

```go
// functional/todo_test.go
package functional

import (
	"net/http"
	"strings"
	"testing"
	
	"../handler"
)

func TestGetSamples(t *testing.T) {
	testServer := setupServer()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/samples", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got := strings.TrimSpace(rec.Body.String())

	want := `[{"id":1,"title":"Do dishes","note":"","due_date":"2000-01-01T00:00:00Z"},{"id":2,"title":"Do homework","note":"","due_date":"2000-01-01T00:00:00Z"},{"id":2,"title":"Twitter","note":"","due_date":"2000-01-01T00:00:00Z"}]`

	if got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}

func setupServer() *http.ServeMux {
	return handler.SetUpRouting()
}
```

This test is very simple. It just calls the endpoint and compares the response body with the expected value. It should pass.

```sh
$ go test ./functional/todo_test.go
ok      command-line-arguments  0.099s
```

## main.go

It's finally time to make `main.go` that call `handler.SetUpRouting()` and set up the HTTP server.

```sh
$ touch main.go
```

```go
// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"./handler"
)

func main() {
	mux := handler.SetUpRouting()

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Run `main.go` and visit `http://localhost:8080/samples`. The sample Json file will be returned.

```json
$ go run main.go
$ curl "http://localhost:8080/samples" | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   221  100   221    0     0  27949      0 --:--:-- --:--:-- --:--:-- 31571
[
  {
    "id": 1,
    "title": "Do dishes",
    "note": "",
    "due_date": "2000-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "title": "Do homework",
    "note": "",
    "due_date": "2000-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "title": "Twitter",
    "note": "",
    "due_date": "2000-01-01T00:00:00Z"
  }
]
```

## PostgreSQL

We've created API that returns TODO list. But it can only return a static file so we can't add a new task, update an existed task, and delete them. We have to implement `Postgres` repository and new endpoints to call its method.


### DB for tests

```sh
$ mkdir testdb
$ touch testdb/testdb.go
```

```go
// testdb/testdb.go
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
```

### Test for Postgres

Create `db/postgres.go` and add methods to implement the `Repository` interface.

```sh
$ touch db/postgres.go
```

```go
// db/postgres.go
package db

import (
	"database/sql"

	"github.com/cohhei/go-to-the-handson/04/schema"
)

type Postgres struct {
	DB *sql.DB
}

func (p *Postgres) Close() {}

func (p *Postgres) Insert(todo *schema.Todo) (int, error) {
	return 0, nil
}

func (p *Postgres) Delete(id int) error {
	return nil
}

func (p *Postgres) GetAll() ([]schema.Todo, error) {
	return nil, nil
}
```

Then create `db/postgres_test.go`.

```sh
$ touch db/postgres_test.go
```

```go
// db/postgres_test.go
package db

import (
	"reflect"
	"testing"
	"time"

	"../schema"
	"../testdb"
	_ "github.com/lib/pq"
)

func TestPostgres_Insert(t *testing.T) {
	postgres := &Postgres{testdb.Setup()}
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
	postgres := &Postgres{testdb.Setup()}
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
	postgres := &Postgres{testdb.Setup()}
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

func equal(got interface{}, want interface{}) bool {
	return reflect.DeepEqual(got, want)
}
```

The test functions contain other methods in the `Postgres` structure for instance the `TestPostgres_Delete` has not only the `Postgres.Delete` method but `Insert` and `GetAll`. That's too bad pattern. When `Insert` has any bugs, `TestPostgres_Delete` will fail even if `Delete` has been correctly implemented. You should build functions to input and output to DB if you can.

Of course, the tests will fail yet.

```sh
$ go test ./db/postgres*
--- FAIL: TestPostgres_Insert (0.03s)
        postgres_test.go:31: <nil>
FAIL
FAIL    command-line-arguments  0.077s
```

### Implementation for Postgres

```go
// db/postgres.go
package db

import (
	"database/sql"

	"../schema"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
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
```

To pass the tests, you should set PostgreSQL up. If you've already installed Docker, you can easily build it by using the one-line command.

```sh
$ docker run -d --name docker-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:alpine
$ go test ./db/postgres*
ok      command-line-arguments  0.250s
```

### Functional test for Postgres

```go
// functional/todo_test.go

// ...
import (
	// ...
	"bytes"
	"fmt"
	"net/http/httptest"
	"reflect"
	"time"
	// ...
	"../db"
	"../schema"
	"../testdb"
	// ...
)
// ...
func TestGetSamples(t *testing.T) {
	testServer := setupServer(nil)
//...
func setupServer(postgres *db.Postgres) *http.ServeMux {
	return handler.SetUpRouting(postgres)
}
// ...
func TestGetAllTodo(t *testing.T) {
	postgres := &db.Postgres{testdb.Setup()}
	testServer := setupServer(postgres)

	todo := &schema.Todo{
		Title:   "My Task1",
		DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
	}

	_, err := postgres.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got := strings.TrimSpace(rec.Body.String())

	want := `[{"id":1,"title":"My Task1","note":"","due_date":"2000-01-01T00:00:00+09:00"}]`

	if got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}

func TestSaveTodo(t *testing.T) {
	postgres := &db.Postgres{testdb.Setup()}
	testServer := setupServer(postgres)

	body := []byte(`{"id":1,"title":"My Task1","note":"","due_date":"2000-01-01T00:00:00+09:00"}`)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/todo", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got := strings.TrimSpace(rec.Body.String())
	want := "1"

	if got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}

	gotTodo, err := postgres.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	wantTodo := []schema.Todo{
		{
			Title:   "My Task1",
			DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want: %v, Got: %v\n", wantTodo, gotTodo)
	}
}

func TestDeleteTodo(t *testing.T) {
	postgres := &db.Postgres{testdb.Setup()}
	testServer := setupServer(postgres)

	todo := &schema.Todo{
		Title:   "My Task1",
		DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
	}

	id, err := postgres.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	body := []byte(fmt.Sprintf(`{"id":%d}`, id))

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:9999/todo", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got := rec.Body.String()

	want := ""

	if got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}

	gotTodo, err := postgres.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(gotTodo) > 0 {
		t.Fatalf("Should return the empty slice, Got: %v\n", gotTodo)
	}
}
```

```sh
$ go test ./functional/todo_test.go
# command-line-arguments
functional/todo_test.go:164:22: too many arguments in call to handler.SetUpRouting
        have (*db.Postgres)
        want ()
FAIL    command-line-arguments [build failed]
```

### New routing and handlers

```go
// handler/todo.go
import (
	// ...
	"io/ioutil"
	// ...
	"../schema"
)

type todoHandler struct {
	postgres *db.Postgres
	samples  *db.Sample
}
//...
func (handler *todoHandler) saveTodo(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetRepository(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var todo schema.Todo
	if err := json.Unmarshal(b, &todo); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := service.Insert(ctx, &todo)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseOk(w, id)
}

func (handler *todoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetRepository(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(b, &req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := service.Delete(ctx, req.ID); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *todoHandler) getAllTodo(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetRepository(r.Context(), handler.postgres)

	todoList, err := service.GetAll(ctx)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseOk(w, todoList)
}
```

We are setting a repository at each handler function, but that's just in order to explain interfaces in Go. So you shouldn't use the pattern in your production code.

```go
// handler/routes.go
package handler

import (
	"net/http"

	"../db"
)

// - func SetUpRouting() {
func SetUpRouting(postgres *db.Postgres) *http.ServeMux {
	todoHandler := &todoHandler{
		postgres: postgres,
		samples:  &db.Sample{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/samples", todoHandler.GetSamples)
	mux.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoHandler.getAllTodo(w, r)
		case http.MethodPost:
			todoHandler.saveTodo(w, r)
		case http.MethodDelete:
			todoHandler.deleteTodo(w, r)
		default:
			responseError(w, http.StatusNotFound, "")
		}
	})

	return mux
}
```

```sh
$ go test ./functional/todo_test.go
ok      command-line-arguments  0.179s
```

After passing functional tests, you can stop the docker container.

```sh
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
530507a60b58        postgres:alpine     "docker-entrypoint.s…"   19 hours ago        Up 19 hours         0.0.0.0:5432->5432/tcp   docker-postgres

$ docker stop 530507a60b58
```

### Updating main.go

```go
// main.go

import (
	//...
	"time"
	"./db"
	//...
)

func main() {
	var postgres *db.Postgres
	var err error
	for i := 0; i < 10; i++ {
		time.Sleep(3 * time.Second)
		postgres, err = db.ConnectPostgres()
	}
	if err != nil {
		panic(err)
	} else if postgres == nil {
		panic("postgres is nil")
	}

	mux := handler.SetUpRouting(postgres)

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

We don't create the `db/ConnectPostgres` function. Add it to `db/postgres.go`.

```go
// db/postgres.go
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
```

What's the difference between the function and `testdb.connectPostgresForTests`? It's the `connStr` string. `localhost:5432` was configured as a host name in `testdb.connectPostgresForTests`. But `postgres` is configured here. Docker will resolve the name.

## Docker

### Postgres

```sh
$ mkdir postgres
$ touch postgres/up.sql
```

```sql
-- postgres/up.sql
DROP TABLE IF EXISTS todo;
CREATE SEQUENCE todo_id START 1;
CREATE TABLE todo (
  ID serial PRIMARY KEY,
  TITLE TEXT NOT NULL,
  NOTE TEXT,
  DUE_DATE TIMESTAMP WITH TIME ZONE
);
```

```sh
$ touch postgres/Dockerfile
```

```dockerfile
# postgres/Dockerfile
FROM postgres:10.3
COPY up.sql /docker-entrypoint-initdb.d/1.sql
CMD ["postgres"]
```

### TODO API

```sh
$ touch Dockerfile
```

```dockerfile
# Dockerfile
# build stage
FROM golang:1.10.2-alpine AS build
ARG dir=/todo
ADD . ${dir}
RUN apk update && \
    apk add --virtual build-dependencies build-base git && \
    cd ${dir} && \
    go get -u github.com/lib/pq && \
    go build -o todo-api

# final stage
FROM alpine:3.7
ARG dir=/todo
WORKDIR /app
COPY --from=build ${dir}/todo-api /app/
EXPOSE 8080
CMD ./todo-api
```

### docker-compose.yml

```sh
$ touch docker-compose.yml
```

```yaml
# docker-compose.yml
version: "3.6"

services:
  postgres:
    build: "./postgres"
    restart: "always"
    environment:
      POSTGRES_DB: "postgres"
      POSTGRES_USER: "postgres"
      POSTGRES_PASSWORD: "postgres"
  todo-api:
    build: "."
    depends_on:
      - postgres
    ports:
      - 8080:8080
```

```sh
$ docker-compose up -d
Creating todo_postgres_1 ... done
Creating todo_todo-api_1 ... done

 $ docker-compose ps
     Name                    Command              State           Ports
--------------------------------------------------------------------------------
todo_postgres_1   docker-entrypoint.sh postgres   Up      5432/tcp
todo_todo-api_1   /bin/sh -c ./todo-api           Up      0.0.0.0:8080->8080/tcp
```

## Requests

We've completed the TODO-API but it can't be used from the browser. Create `requests/request.go`, a cli tool to use our API.

```sh
$ mkdir requests
$ touch requests/requests.go
```

```go
// requests/requests.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const usage = `
usage: todo COMMAND

Commands:
  samples     Get sample todo tasks
  all         Get all todo tasks
  add         Add new todo task
  delete      Remove a todo task
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		return
	}

	switch command := os.Args[1]; command {
	case "samples":
		get("samples")
	case "all":
		get("todo")
	case "add":
		add()
	case "delete":
		del()
	default:
		fmt.Printf("'%s' is not a todo command.", command)
	}
}

func get(path string) {
	res, err := http.Get(fmt.Sprintf("http://localhost:8080/%s", path))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if len(b) == 0 {
		fmt.Println("The response is empty.")
		return
	}
	fmt.Println(string(b))
}

const usage_add = `
usage: todo add TODO_NAME TODO_NOTE DUE_DATE
`

func add() {
	if len(os.Args) < 3 {
		fmt.Print(usage_add)
		return
	}

	name := os.Args[2]
	var note string
	var date string

	if len(os.Args) > 3 {
		note = os.Args[3]
		if len(os.Args) > 4 {
			date = os.Args[4]
		}
	}

	todo := struct {
		Title   string `json:"title"`
		Note    string `json:"note,omitempty"`
		DueDate string `json:"due_date,omitempty"`
	}{
		name, note, date,
	}

	b, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	res, err := http.Post("http://localhost:8080/todo", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

func del() {
	if len(os.Args) < 3 {
		fmt.Print("usage: todo delete TASK_ID")
		return
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("TASK_ID should be number")
		return
	}

	b, err := json.Marshal(map[string]int{"id": id})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todo", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	res.Body.Close()
}
```

You can call all API's endpoints by sub commands, `samples`, `all`, `add`, and `delete`. 

```sh
$ go run requests/requests.go

usage: todo COMMAND

Commands:
  samples     Get sample todo tasks
  all         Get all todo tasks
  add         Add new todo task
  delete      Remove a todo task
```

For instance, add a new task by `add $TASK_NAME $NOTE` and get all tasks by `all`.

```json
$ go run requests/requests.go add Task1 'this is the first task.'
$ go run requests/requests.go all | jq
[
  {
    "id": 1,
    "title": "Task1",
    "note": "this is the first task.",
    "due_date": "2000-01-01T00:00:00Z"
  }
]
```

[PREV](../03)