package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cohhei/go-to-the-handson/04/db"
	"github.com/cohhei/go-to-the-handson/04/handler"
)

func main() {
	var postgres *db.Postgres
	var err error
	for i := 0; i < 10; i++ {
		time.Sleep(3 * time.Second)
		postgres, err = db.ConnectPostgres()
	}
	if err != nil {
		panic(err)
	} else if postgres == nil {
		panic("postgres is nil")
	}

	mux := handler.SetUpRouting(postgres)

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
