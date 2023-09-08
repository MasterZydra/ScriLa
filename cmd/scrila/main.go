package main

import (
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// fmt.Println(lexer.Tokenize("int x = 42;"))
	repl()
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	parser := parser.NewParser()

	fmt.Println("\nRepl v0.1")
	for true {
		input := readInput(reader)

		if input == "" || input == "exit" {
			os.Exit(0)
		}

		fmt.Println("Tokens:", lexer.Tokenize(input))

		program := parser.ProduceAST(input)
		fmt.Printf("AST: [%s, Body: %s\n", program.GetKind(), program.GetBody())
	}
}

func readInput(reader *bufio.Reader) string {
	fmt.Print("\n> ")
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}
