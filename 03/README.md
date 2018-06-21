# 3. Connect DB and Execute queries

## setup

```sh
$ docker run -d --name docker-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:alpine
...

$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
d9065d8c1f30        postgres:alpine     "docker-entrypoint.s…"   3 seconds ago       Up 4 seconds        0.0.0.0:5432->5432/tcp   docker-postgres
```

## Connect PostgreSQL

```sh
$ touch connect.go
```

```go
// connect.go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	// connStr := "user=postgres dbname=postgres sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	
	// [user]:[password]@[address]/[DB name]
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Ping OK")
}
```

```sh
$ go run connectdb.go
Ping OK
```

## Create Tables

```sh
$ touch queries.go
```

```go
// queries.go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const createTable = `
DROP TABLE IF EXISTS ACCOUNT;
CREATE TABLE ACCOUNT
(
	ID serial,
	NAME varchar(50),
	MAIL_ADDRESS varchar(50),
	LANG varchar(5)
);
`

func main() {
	db, err := connectPostgres()
	if err != nil {
		return
	}

	defer db.Close()

	if _, err = db.Exec(createTable); err != nil {
		fmt.Println(err)
		return
	}

	if _, err = db.Exec("DROP TABLE ACCOUNT"); err != nil {
		fmt.Println(err)
		return
	}
}

func connectPostgres() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Ping OK")

	return db, nil
}
```

```sh
$ go run queries.go
Ping OK
```

## Exercise 3-1

Add a function `insertAccount()` into `queries.go`. It inserts records which contain data of any static accounts to the `account` table by inside a transaction.

```go
// queries.go
func main() {
	db, err := connectPostgres()
	if err != nil {
		return
	}

	defer db.Close()

	if _, err = db.Exec(createTable); err != nil {
		fmt.Println(err)
		return
	}

	if err := insertAccounts(db); err != nil {
		fmt.Println(err)
		return
	}

	if _, err = db.Exec("DROP TABLE ACCOUNT"); err != nil {
		fmt.Println(err)
		return
	}
}

// You should implement the function
func insertAccounts(db *sql.DB) error {}
```

HINTS: You can create a transaction with [`DB.begin`](https://golang.org/pkg/database/sql/#DB.Begin) and commit it with [`Tx.Commit`](https://golang.org/pkg/database/sql/#Tx.Commit).

## Exercise 3-2

Add a function `getAccounts()` into `queries.go`. It returns all records in the `account` table.

```go
// queries.go
	if err := insertAccounts(db); err != nil {
		fmt.Println(err)
		return
	}

	accounts, err := getAccounts(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(accounts)
...

// You should implement the function and define the Account struct
func getAccounts(db *sql.DB) ([]Account, error) {}
```

```sh
$ go run queries.go
Ping OK
[{1 My Name my_name@example.com ja} {2 Your Name your_name@example.com en}]
```

HINTS: You can execute `SELECT` queries with `DB.Query` and get the results with [`Rows.Scan`](https://golang.org/pkg/database/sql/#Row.Scan).

The answer is [queries.go](main/queries.go).

[PREV](../02) | [NEXT](../04)
