package main

import (
	"flag"
	"fmt"
)

func main() {
	// You can get command line options by flag package.
	// 'flag.StringVar' returns a string option as a pointer.
	// If you want to know other flag package's functions, go to https://golang.org/pkg/flag
	var name string
	flag.StringVar(&name, "name", "", "Write your name.")

	flag.Parse()

	if name == "" {
		fmt.Println("Hello World!")
	} else {
		fmt.Printf("Hello %s!\n", name)
	}
}
