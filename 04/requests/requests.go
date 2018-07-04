package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const usage = `
usage: todo COMMAND

Commands:
  samples     Get sample todo tasks
  all         Get all todo tasks
  add         Add new todo task
  delete      Remove a todo task
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		return
	}

	switch command := os.Args[1]; command {
	case "samples":
		get("samples")
	case "all":
		get("todo")
	case "add":
		add()
	case "delete":
		del()
	default:
		fmt.Printf("'%s' is not a todo command.", command)
	}
}

func get(path string) {
	res, err := http.Get(fmt.Sprintf("http://localhost:8080/%s", path))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if len(b) == 0 {
		fmt.Println("The response is empty.")
		return
	}
	fmt.Println(string(b))
}

const usage_add = `
usage: todo add TODO_NAME TODO_NOTE DUE_DATE
`

func add() {
	if len(os.Args) < 3 {
		fmt.Print(usage_add)
		return
	}

	name := os.Args[2]
	var note string
	var date string

	if len(os.Args) > 3 {
		note = os.Args[3]
		if len(os.Args) > 4 {
			date = os.Args[4]
		}
	}

	todo := struct {
		Title   string `json:"title"`
		Note    string `json:"note,omitempty"`
		DueDate string `json:"due_date,omitempty"`
	}{
		name, note, date,
	}

	b, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	res, err := http.Post("http://localhost:8080/todo", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

func del() {
	if len(os.Args) < 3 {
		fmt.Print("usage: todo delete TASK_ID")
		return
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("TASK_ID should be number")
		return
	}

	b, err := json.Marshal(map[string]int{"id": id})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todo", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	res.Body.Close()
}
