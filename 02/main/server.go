package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			MailAddress string `json:"mail_address"`
			Options     string `json:"options,omitempty"`
		}{
			ID:   1,
			Name: "my_name",
		}

		b, err := json.Marshal(data)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))

		w.Write(b)
	})

	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found, %s", r.URL)
	})

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
