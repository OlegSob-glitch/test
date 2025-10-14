package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "World", "Имя для приветствия")
	flag.Parse()
	fmt.Printf("Hello, %s! from first_app\n", *name)
}
