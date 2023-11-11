package main

import (
	"ScriLa/cmd/scrila/bashTranspiler"
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/parser"
	"flag"
	"fmt"
	"os"
)

func main() {
	showTokens := flag.Bool("st", false, "Show tokens")
	showAST := flag.Bool("sa", false, "Show AST")
	showCallStack := flag.Bool("sc", false, "Show call stack")
	filename := flag.String("f", "", "Script file")
	flag.Parse()

	// Check if filename is not empty
	if *filename == "" {
		fmt.Println("Usage: scrila -f [filename.scri]")
		os.Exit(1)
	}

	transpile(*filename, *showTokens, *showAST, *showCallStack)
}

func transpile(filename string, showTokens bool, showAST bool, showCallStack bool) {
	parser := parser.NewParser()
	transpilerObj := bashTranspiler.NewTranspiler(showCallStack)
	env := bashTranspiler.NewEnvironment(nil, transpilerObj)

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
