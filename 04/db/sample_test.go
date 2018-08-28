package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/cohhei/go-to-the-handson/04/schema"
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
