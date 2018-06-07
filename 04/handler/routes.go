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
}
