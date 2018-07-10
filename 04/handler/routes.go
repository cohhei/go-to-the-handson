package handler

import (
	"net/http"

	"github.com/cohhei/go-to-the-handson/04/db"
)

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
