package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	name := flag.String("name", "World", "Имя для приветствия")
	flag.Parse()
	fmt.Printf("Hello, %s! from first_app\n", *name)

	// Вывод текущего времени
	fmt.Println("Current time:", time.Now())
}
