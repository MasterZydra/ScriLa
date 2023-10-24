package main

import (
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/parser"
	"ScriLa/cmd/scrila/transpiler"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Check if a argument is passed
	if len(os.Args) != 2 {
		fmt.Println("Usage: scrila [filename.scri]")
		os.Exit(1)
	}
	filename := os.Args[1]
	// Check if filename is not empty
	if filename == "" {
		fmt.Println("Usage: scrila [filename.scri]")
		os.Exit(1)
	}

	showTokens := flag.Bool("st", false, "Show tokens")
	showAST := flag.Bool("sa", false, "Show AST")
	flag.Parse()

	transpile(filename, *showTokens, *showAST)
}

func transpile(filename string, showTokens bool, showAST bool) {
	parser := parser.NewParser()
	transpilerObj := transpiler.NewTranspiler()
	env := transpiler.NewEnvironment(nil, transpilerObj)

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file '" + filename + "':")
		fmt.Println(err)
		os.Exit(1)
	}

	if showTokens {
		tokens, err := lexer.NewLexer().Tokenize(string(fileContent), filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Tokens:   %s\n", tokens)
	}
	program, err := parser.ProduceAST(string(fileContent), filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if showAST {
		fmt.Printf("AST:       %s\n", program)
	}
	err = transpilerObj.Transpile(program, env, filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
