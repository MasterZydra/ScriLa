package main

import (
	"ScriLa/cmd/scrila/lexer"
	"fmt"
)

func main() {
	fmt.Println(lexer.Tokenize("int x = 42;"))
}
