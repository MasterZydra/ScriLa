package main

import (
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/parser"
	"ScriLa/cmd/scrila/runtime"
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
	lexer := lexer.NewLexer()
	parser := parser.NewParser()
	env := runtime.NewEnvironment(nil)

	fmt.Println("\nRepl v0.1")
	for true {
		input := readInput(reader)

		if input == "" || input == "exit" {
			os.Exit(0)
		}

		fmt.Printf("Tokens:   %#v\n", lexer.Tokenize(input))

		program := parser.ProduceAST(input)
		fmt.Printf("AST:       %s\n", program)

		result := runtime.Evaluate(program, env)
		fmt.Println("Interpret:", result)
	}
}

func readInput(reader *bufio.Reader) string {
	fmt.Print("\n> ")
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}
