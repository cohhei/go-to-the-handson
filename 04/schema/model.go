package schema

import "time"

type Todo struct {
	ID      int
	Title   string
	Note    string
	DueDate time.Time
}
