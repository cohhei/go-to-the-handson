package db

import (
	"time"

	"github.com/cohhei/go-to-the-handson/04/schema"
)

type Sample struct{}

func (s *Sample) Close() {}

func (s *Sample) Insert(todo *schema.Todo) (int, error) {
	return 0, nil
}

func (s *Sample) Delete(id int) error {
	return nil
}

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
