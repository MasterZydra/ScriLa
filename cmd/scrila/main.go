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
	fileName := flag.String("f", "", "Path of file")
	showTokens := flag.Bool("st", false, "Show tokens")
	showAST := flag.Bool("sa", false, "Show AST")
	flag.Parse()

	if *fileName == "" {
		fmt.Println("Usage of scrila:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	transpile(*fileName, *showTokens, *showAST)
}

func transpile(filename string, showTokens bool, showAST bool) {
	parser := parser.NewParser()
	env := transpiler.NewEnvironment(nil)

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
	err = transpiler.Transpile(program, env, filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
