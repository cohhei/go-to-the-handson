package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cohhei/go-to-the-handson/04/db"
	"github.com/cohhei/go-to-the-handson/04/schema"
	"github.com/cohhei/go-to-the-handson/04/service"
)

type todoHandler struct {
	postgres *db.Postgres
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
