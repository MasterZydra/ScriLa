package lexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Lexer struct {
	sourceChars []string
	tokens      []*Token
	currLn      int
	currCol     int
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (self *Lexer) init() {
	self.tokens = make([]*Token, 0)
	self.currLn = 1
	self.currCol = 1
}

func (self *Lexer) Tokenize(sourceCode string) []*Token {
	self.init()

	// Split source code into an array of every character
	self.sourceChars = strings.Split(sourceCode, "")

	for self.isNotEof() {
		// Handle single-character tokens
		if reserved, ok := singleCharTokens[self.at()]; ok {
			self.pushToken(self.eat(), reserved)
			continue
		}

		// --- Handle multi-character tokens ---

		// Build Int token
		if isDigit(self.at()) {
			num := ""
			for self.isNotEof() && isDigit(self.at()) {
				num += self.eat()
			}
			self.pushToken(num, Int)
			continue
		}

		if isLetter(self.at()) {
			ident := ""
			for self.isNotEof() && isLetter(self.at()) {
				ident += self.eat()
			}

			// Check for reserved keywords
			if reserved, ok := keywords[ident]; ok {
				self.pushToken(ident, reserved)
			} else {
				self.pushToken(ident, Identifier)
			}
			continue
		}

		if isSkippable(self.at()) {
			if self.at() == "\n" {
				self.currLn += 1
				self.currCol = 0
			}

			self.eat()
			continue
		}

		fmt.Println("Unrecognized character found:", self.at())
		os.Exit(1)
	}

	self.pushToken("EOF", EndOfFile)

	return self.tokens
}

func (self *Lexer) isNotEof() bool {
	return len(self.sourceChars) > 0
}

func (self *Lexer) at() string {
	return self.sourceChars[0]
}

func (self *Lexer) eat() string {
	self.currCol += 1
	var prev string
	prev, self.sourceChars = self.sourceChars[0], self.sourceChars[1:]
	return prev
}

func (self *Lexer) pushToken(value string, tokenType TokenType) {
	col := self.currCol
	if tokenType != EndOfFile {
		col -= len(value)
	}
	self.tokens = append(self.tokens, &Token{
		TokenType: tokenType,
		Value:     value,
		Ln:        self.currLn,
		Col:       col,
	})
}

func isLetter(sourceChar string) bool {
	return unicode.IsLetter([]rune(sourceChar)[0])
}

func isDigit(sourceChar string) bool {
	return unicode.IsDigit([]rune(sourceChar)[0])
}

func isSkippable(sourceChar string) bool {
	return sourceChar == " " || sourceChar == "\r" || sourceChar == "\n" || sourceChar == "\t"
}
