package lexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type TokenType string

const (
	Semicolon TokenType = "Semicolon"
	EndOfFile TokenType = "EOF"
	// --- Operations ---
	BinaryOperator TokenType = "BinaryOperator"
	// --- Priority ---
	OpenParen  TokenType = "OpenParen"
	CloseParen TokenType = "CloseParen"
	// --- Variables ---
	Identifier TokenType = "Identifier"
	Equals     TokenType = "Equals"
	// Variable types
	Bool     TokenType = "BoolValue"
	BoolType TokenType = "BoolType"
	Int      TokenType = "IntValue"
	IntType  TokenType = "IntType"
	NullType TokenType = "NullType"
	Str      TokenType = "StrValue"
	StrType  TokenType = "StrType"
)

var singleCharTokens = map[string]TokenType{
	"-": BinaryOperator,
	";": Semicolon,
	"(": OpenParen,
	")": CloseParen,
	"*": BinaryOperator,
	"/": BinaryOperator,
	"+": BinaryOperator,
	"=": Equals,
}

var keywords = map[string]TokenType{
	"bool": BoolType,
	"int":  IntType,
	"null": NullType,
	"str":  StrType,
}

type Token struct {
	TokenType TokenType
	Value     string
}

func Tokenize(sourceCode string) []Token {
	tokens := make([]Token, 0)

	// Split source code into an array of every character
	sourceChars := strings.Split(sourceCode, "")

	// Build each token until EOF
	for len(sourceChars) > 0 {
		// --- Handle single-character tokens ---
		if reserved, ok := singleCharTokens[sourceChars[0]]; ok {
			tokens = append(tokens, createToken(sourceChars[0], reserved))
			sourceChars = removeFirstElem(sourceChars)
			continue
		}

		// --- Handle multi-character tokens ---

		// Build Int token
		if isDigit(sourceChars[0]) {
			num := ""
			for len(sourceChars) > 0 && isDigit(sourceChars[0]) {
				num += sourceChars[0]
				sourceChars = removeFirstElem(sourceChars)
			}
			tokens = append(tokens, createToken(num, Int))
			continue
		}

		if isLetter(sourceChars[0]) {
			ident := ""
			for len(sourceChars) > 0 && isLetter(sourceChars[0]) {
				ident += sourceChars[0]
				sourceChars = removeFirstElem(sourceChars)
			}

			// Check for reserved keywords
			if reserved, ok := keywords[ident]; ok {
				tokens = append(tokens, createToken(ident, reserved))
			} else {
				tokens = append(tokens, createToken(ident, Identifier))
			}
			continue
		}

		if isSkippable(sourceChars[0]) {
			sourceChars = removeFirstElem(sourceChars)
			continue
		}

		fmt.Println("Unrecognized character found:", sourceChars[0])
		os.Exit(1)
	}

	tokens = append(tokens, createToken("EOF", EndOfFile))

	return tokens
}

func createToken(value string, tokenType TokenType) Token {
	return Token{
		TokenType: tokenType,
		Value:     value,
	}
}

func removeFirstElem(sourceChars []string) []string {
	return sourceChars[1:]
}

func isLetter(sourceChar string) bool {
	return unicode.IsLetter([]rune(sourceChar)[0])
}

func isDigit(sourceChar string) bool {
	return unicode.IsDigit([]rune(sourceChar)[0])
}

func isSkippable(sourceChar string) bool {
	return sourceChar == " " || sourceChar == "\n" || sourceChar == "\t"
}

// int x = 42;
// int y = (1 + 2);
