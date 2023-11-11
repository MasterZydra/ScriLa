package main

import (
	"ScriLa/cmd/scrila/bashAssembler"
	"ScriLa/cmd/scrila/bashTranspiler"
	"ScriLa/cmd/scrila/config"
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

	config.Filename = *filename

	transpile(*showTokens, *showAST, *showCallStack)
}

func transpile(showTokens bool, showAST bool, showCallStack bool) {
	parser := parser.NewParser()
	transpilerObj := bashTranspiler.NewTranspiler(showCallStack)
	env := bashTranspiler.NewEnvironment(nil, transpilerObj)

	fileContent, err := os.ReadFile(config.Filename)
	if err != nil {
		fmt.Println("Error reading file '" + config.Filename + "':")
		fmt.Println(err)
		os.Exit(1)
	}

	if showTokens {
		tokens, err := lexer.NewLexer().Tokenize(string(fileContent))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Tokens:   %s\n", tokens)
	}
	program, err := parser.ProduceAST(string(fileContent))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if showAST {
		fmt.Printf("AST:       %s\n", program)
	}
	bashProgram, err := transpilerObj.Transpile(program, env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	assembler := bashAssembler.NewAssembler()
	err = assembler.Assemble(bashProgram)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
