package main

import (
	"flag"
	"fmt"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "Write your name.")

	flag.Parse()

	if name == "" {
		fmt.Println("Hello World!")
	} else {
		fmt.Printf("Hello %s!\n", name)
	}
}
