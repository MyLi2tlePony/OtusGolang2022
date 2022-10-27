package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Количество аргументов должно быть больше 1")
		return
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
