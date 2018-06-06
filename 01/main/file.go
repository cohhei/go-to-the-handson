package main

import (
	"fmt"
	"os"
)

func main() {
	// You can get command line arguments with 'os.Args', a string slice.
	if len(os.Args) < 2 {
		fmt.Println("Please set a file name.")
		return
	}

	// 'os.Args' contains the executed binary file name and the arguments.
	// If you command './file file.txt', 'os.Args[0]' is './file' and 'os.Args[1] is 'file.txt'.
	filename := os.Args[1]

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	b := []byte("Hello Writing to Files!")
	n, err := f.Write(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("The number of bytes written: ", n)
}
