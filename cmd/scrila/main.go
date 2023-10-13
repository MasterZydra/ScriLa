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
	filename := flag.String("f", "", "Path of file")
	showTokens := flag.Bool("st", false, "Show tokens")
	showAST := flag.Bool("sa", false, "Show AST")
	flag.Parse()

	if *filename == "" {
		fmt.Println("Usage of scrila:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	transpile(*filename, *showTokens, *showAST)
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
