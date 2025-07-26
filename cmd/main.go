package main

import (
	"fmt"
	"os"
)

func main() {
	dir := os.Args[1]
	fmt.Println("Setting up project directory:", dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			panic(err)
		}
	}
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}


	Run()
}
