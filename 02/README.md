# 2. Create a HTTP request and server

## `nat/http` package

### Create a request

```go
// request.go
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

	fmt.Println(string(body))
}
```

```json
$ go run request.go | jq
{
  "login": "defunkt",
  "id": 2,
  "node_id": "MDQ6VXNlcjI=",
  "avatar_url": "https://avatars0.githubusercontent.com/u/2?v=4",
  "gravatar_id": "",
  "url": "https://api.github.com/users/defunkt",
  "html_url": "https://github.com/defunkt",
  "followers_url": "https://api.github.com/users/defunkt/followers",
  "following_url": "https://api.github.com/users/defunkt/following{/other_user}",
  "gists_url": "https://api.github.com/users/defunkt/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/defunkt/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/defunkt/subscriptions",
  "organizations_url": "https://api.github.com/users/defunkt/orgs",
  "repos_url": "https://api.github.com/users/defunkt/repos",
  "events_url": "https://api.github.com/users/defunkt/events{/privacy}",
  "received_events_url": "https://api.github.com/users/defunkt/received_events",
  "type": "User",
  "site_admin": true,
  "name": "Chris Wanstrath",
  "company": "@github ",
  "blog": "http://chriswanstrath.com/",
  "location": "San Francisco",
  "email": null,
  "hireable": true,
  "bio": "🍔 ",
  "public_repos": 107,
  "public_gists": 273,
  "followers": 20213,
  "following": 210,
  "created_at": "2007-10-20T05:24:19Z",
  "updated_at": "2018-06-05T19:29:51Z"
}
```

### Setup a server

```go
// server.go
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

```sh
$ go run server.go
http://localhost:8080
...

$ curl http://localhost:8080
Hello, "/"
```

## `json` package

```go
// json.go
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

```sh
$ go run json.go
ID:  1 Name:  Gopher
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