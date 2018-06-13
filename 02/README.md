# 2. Create a HTTP request and server

## `nat/http` package

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	res, err := http.Get("https://api.github.com/users/defunkt")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(body)
}
```

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## `json` package

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	str := []byte(`{"id":1,"name":"Gopher"}`)
	data := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{}

	if err := json.Unmarshal(str, &data); err != nil {
		panic(data)
	}
	fmt.Println("ID: ", data.ID, "Name: ", data.Name)
}
```

## Exercise 2-1

Create an application `request.go` which creates a `GET` request to https://api.github.com/users/defunkt and parses the response body with `json.Unmarshal` only `login`, `id`, `site_admin`, and `bio`.

```sh
$ go run request.go
Login:  defunkt
ID:  2
SiteAdmin:  true
Bio:  🍔
```

The answer is [request.go](main/request.go)

## Exercise 2-2

Create an application `server.go` which builds a server. It has three endpoints, `/`, `/json`, and `/404`. `/` returns `Hello, "/"`. `/json` returns `{"id":1,"name":"my_name","mail_address":""}`. `/404` returns `Not Found, /404` with the 404 status code.

```sh
$ go run server.go
http://localhost:8080

$ curl http://localhost:8080/
Hello, "/"

$ curl http://localhost:8080/json
{"id":1,"name":"my_name","mail_address":""}

$ curl http://localhost:8080/404
Not Found, /404
```

The answer is [server.go](main/server.go).

[PREV](../01) | [NEXT](../03)