package main

import (
	"ac/internal/role/manager"
	"os"
)

func main() {
	dir := "workdir"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			panic(err)
		}
	}
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}

	manager.A()
}
