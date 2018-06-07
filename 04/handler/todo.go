package handler

import (
	"encoding/json"
	"net/http"

	"github.com/corhhey/go-to-the-handson/04/db"
	"github.com/corhhey/go-to-the-handson/04/service"
)

type todoHandler struct {
	samples *db.Sample
}

func (handler *todoHandler) GetSamples(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetRepository(r.Context(), handler.samples)

	todoList, err := service.GetAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(todoList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(b)
}
