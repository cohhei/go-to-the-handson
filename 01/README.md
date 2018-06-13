# 1. Command Line Arguments and File I/O

## Setup

```sh
$ go get -u github.com/corhhey/go-to-the-handson
...
$ mkdir go-handson
$ cd go-handson
```

## Hello World

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
}
```

```sh
$ go run hello.go
Hello World!

$ go build -o hello .

$ ./hello
Hello World!
```

## `flag` package

### Usage of `flag.StringVar`

```go
package main

import (
	"flag"
	"fmt"
)

func main() {
	var name string
	flag.StringVar(&name, "opt", "", "Usage")

	flag.Parse()

	fmt.Println(name)
}
```

```sh
$ go run -opt option hello.go
option
```

If you want to know more about the `flag` package, please go to the https://golang.org/pkg/flag/

### Exercise 1-1

Create a CLI application which outputs `Hello World!` if no options are specified. And if a string option is specified as `-name`, it has to output `Hello [YOUR_NAME]!`

```sh
$ ./hello
Hello World!

$ ./hello -name Gopher
Hello Gopher!
```

The answer is [hello.go](main/hello.go)

## `os` package
### Usage of `os.Args`

```go
// args.go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args)
	fmt.Println(os.Args[1])
}
```

```sh
$ go build -o args args.go 
$ ./args Gopher
[./args Gopher]
Gopher
```

### File I/O

#### Reading files
```go
file, err := os.Open(`/path/to/file`)
if err != nil {
		panic(err)
}
defer file.Close()

buf := make([]byte, BUFSIZE)
for {
	n, err := file.Read(buf)
	if n == 0 {
		break
	}
	if err != nil {
		panic(err)
	}

	fmt.Print(string(buf[:n]))
}
```

#### Writing files
```go
f, err := os.Create("/path/to/file")
if err != nil {
	panic(err)
}
defer f.Close()

b := []byte("Foo")
n, err := f.Write(b)
if err != nil {
	panic(err)
}
fmt.Println(n)
```

### Exercise 1-2

Create an application `file.go` which creates a file and write a string `Hello Writing to Files!` to it. And the file name has to be specified as a command line argument.

```sh
$ go run file.go file.txt
The number of bytes written:  23

$ cat file.txt
Hello Writing to Files!
```

The answer is [file.go](main/file.go)

[NEXT](../02)