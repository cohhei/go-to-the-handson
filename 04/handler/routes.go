package handler

import (
	"net/http"

	"github.com/corhhey/go-to-the-handson/04/db"
)

func SetUpRouting(postgres *db.Postgres) {
	todoHandler := &todoHandler{
		postgres: postgres,
		samples:  &db.Sample{},
	}

	http.HandleFunc("/samples", todoHandler.GetSamples)
	http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
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
}
