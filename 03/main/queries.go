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

const insertAccount = `
INSERT INTO ACCOUNT (ID, NAME, MAIL_ADDRESS, LANG)
VALUES (1, 'My Name', 'my_name@example.com', 'ja'),
(2, 'Your Name', 'your_name@example.com', 'en');
`

const selectAll = `
SELECT *
FROM ACCOUNT
ORDER BY ID
`

type Account struct {
	ID          int
	Name        string
	MailAddress string `db:"mail_address"`
	Lang        string
}

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

	if _, err = db.Exec(insertAccount); err != nil {
		fmt.Println(err)
		return
	}

	accounts, err := getAccounts(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(accounts)

	if _, err = db.Exec("DROP TABLE ACCOUNT"); err != nil {
		fmt.Println(err)
		return
	}
}

func getAccounts(db *sql.DB) ([]Account, error) {
	rows, err := db.Query(selectAll)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var accounts []Account

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.MailAddress, &account.Lang); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
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
