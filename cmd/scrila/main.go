package main

import (
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/parser"
	"ScriLa/cmd/scrila/runtime"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var printTokens bool
var printAST bool

func main() {
	shallRepl := flag.Bool("repl", false, "Run repl")
	shallInterprete := flag.Bool("i", false, "Shall the input be interpreted")
	fileName := flag.String("f", "", "Path of file")
	showTokens := flag.Bool("st", false, "Show tokens")
	showAST := flag.Bool("sa", false, "Show AST")
	flag.Parse()

	printTokens = *showTokens
	printAST = *showAST

	if *shallRepl {
		repl()
		os.Exit(0)
	}

	if *fileName == "" {
		fmt.Println("Usage of scrila:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *shallInterprete {
		runFile(*fileName)
		os.Exit(0)
	}
}

func runFile(filename string) {
	parser := parser.NewParser()
	env := runtime.NewEnvironment(nil)

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file '" + filename + "':")
		fmt.Println(err)
		os.Exit(1)
	}

	if printTokens {
		fmt.Printf("Tokens:   %s\n", lexer.NewLexer().Tokenize(string(fileContent)))
	}
	program := parser.ProduceAST(string(fileContent))
	if printAST {
		fmt.Printf("AST:       %s\n", program)
	}
	runtime.Evaluate(program, env)
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

		if printTokens {
			fmt.Printf("Tokens:   %s\n", lexer.Tokenize(input))
		}

		program := parser.ProduceAST(input)
		if printAST {
			fmt.Printf("AST:       %s\n", program)
		}

		runtime.Evaluate(program, env)
	}
}

func readInput(reader *bufio.Reader) string {
	fmt.Print("\n> ")
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}
