package main

import (
	"ScriLa/cmd/scrila/bashAssembler"
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/bashTranspiler"
	"ScriLa/cmd/scrila/config"
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/parser"
	"ScriLa/cmd/scrila/scrilaAst"
	"flag"
	"fmt"
	"os"
)

func main() {
	showTokens := flag.Bool("st", false, "Show tokens")
	showAstScriLa := flag.Bool("sas", false, "Show AST for ScriLa")
	showAstBash := flag.Bool("sab", false, "Show AST for Bash")
	showCallStackScrila := flag.Bool("scs", false, "Show call stack for ScriLa")
	showCallStackBash := flag.Bool("scb", false, "Show call stack for Bash")
	filename := flag.String("f", "", "Script file")
	flag.Parse()

	// Check if filename is not empty
	if *filename == "" {
		fmt.Println("Usage: scrila -f [filename.scri]")
		os.Exit(1)
	}

	config.Filename = *filename
	config.ShowTokens = *showTokens
	config.ShowAstScriLa = *showAstScriLa
	config.ShowAstBash = *showAstBash
	config.ShowCallStackScrila = *showCallStackScrila
	config.ShowCallStackBash = *showCallStackBash

	transpile()
}

func transpile() {
	parser := parser.NewParser()
	transpilerObj := bashTranspiler.NewTranspiler()
	env := bashTranspiler.NewEnvironment(nil, transpilerObj)

	fileContent, err := os.ReadFile(config.Filename)
	if err != nil {
		fmt.Println("Error reading file '" + config.Filename + "':")
		fmt.Println(err)
		os.Exit(1)
	}

	if config.ShowTokens {
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
	if config.ShowAstScriLa {
		fmt.Printf("ScriLa-AST:\n%s\n", scrilaAst.SprintAST(program))
	}
	bashProgram, err := transpilerObj.Transpile(program, env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if config.ShowAstBash {
		fmt.Printf("Bash-AST:\n%s\n", bashAst.SprintAst(bashProgram))
	}
	assembler := bashAssembler.NewAssembler()
	err = assembler.Assemble(bashProgram)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
