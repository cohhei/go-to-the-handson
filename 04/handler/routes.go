package handler

import (
	"net/http"

	"github.com/corhhey/go-to-the-handson/04/db"
)

func SetUpRouting() {
	todoHandler := &todoHandler{
		samples: &db.Sample{},
	}

	http.HandleFunc("/samples", todoHandler.GetSamples)
}
