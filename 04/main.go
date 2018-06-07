package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/corhhey/go-to-the-handson/04/db"
	"github.com/corhhey/go-to-the-handson/04/handler"
)

func main() {
	postgres, err := db.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	handler.SetUpRouting(postgres)

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
