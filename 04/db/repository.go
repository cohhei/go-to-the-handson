package db

import (
	"context"

	"github.com/cohhei/go-to-the-handson/04/schema"
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
