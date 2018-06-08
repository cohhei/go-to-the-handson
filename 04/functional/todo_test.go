package functional

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/corhhey/go-to-the-handson/04/handler"
	"github.com/corhhey/go-to-the-handson/04/schema"
	"github.com/corhhey/go-to-the-handson/04/testdb"
)

func TestGetSamples(t *testing.T) {
	testdb.Setup()

	res, err := http.Get("http://localhost:8080/samples")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(b))

	want := `[{"id":1,"title":"Do dishes","note":"","due_date":"2000-01-01T00:00:00Z"},{"id":2,"title":"Do homework","note":"","due_date":"2000-01-01T00:00:00Z"},{"id":2,"title":"Twitter","note":"","due_date":"2000-01-01T00:00:00Z"}]`

	if got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}

func TestGetAllTodo(t *testing.T) {
	testdb.Setup()

	todo := &schema.Todo{
		Title:   "My Task1",
		DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
	}

	_, err := postgres.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://localhost:8080/todo")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(b))

	want := `[{"id":1,"title":"My Task1","note":"","due_date":"2000-01-01T00:00:00+09:00"}]`

	if got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}

func TestSaveTodo(t *testing.T) {
	testdb.Setup()

	body := []byte(`{"id":1,"title":"My Task1","note":"","due_date":"2000-01-01T00:00:00+09:00"}`)

	res, err := http.Post("http://localhost:8080/todo", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(b))

	want := `1`

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
	testdb.Setup()

	todo := &schema.Todo{
		Title:   "My Task1",
		DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
	}

	id, err := postgres.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	body := []byte(fmt.Sprintf(`{"id":%d}`, id))

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todo", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(string(b))

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

var postgres = testdb.Setup()

var _ = setupServer()

func setupServer() error {
	handler.SetUpRouting(postgres)

	go http.ListenAndServe(":8080", nil)

	return nil
}
