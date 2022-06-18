package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	phrase := "Hello, OTUS!"
	reversePhrase := stringutil.Reverse(phrase)
	fmt.Print(reversePhrase)
}
