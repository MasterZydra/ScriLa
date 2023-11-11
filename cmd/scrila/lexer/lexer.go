package lexer

import (
	"ScriLa/cmd/scrila/config"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
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

func (self *Lexer) Tokenize(sourceCode string) ([]*Token, error) {
	self.init()

	// Split source code into an array of every character
	self.sourceChars = strings.Split(sourceCode, "")

	for self.isNotEof() {
		if self.at() == "#" {
			self.tokenizeComment()
			continue
		}

		if self.at() == "\"" {
			self.tokenizeString()
			continue
		}

		// Handle boolean and comparison operations
		if op := self.at() + self.next(0); slices.Contains(BooleanOps, op) || slices.Contains(ComparisonOps, op) {
			operation := self.eat() + self.eat()
			self.pushToken(operation, BinaryOperator)
		}

		// Handle single-character tokens
		if reserved, ok := singleCharTokens[self.at()]; ok {
			currChar := self.eat()

			// Resolve short form operators like +=, -=, *=, /=, ...
			if reserved == BinaryOperator && self.at() == "=" {
				self.eat()
				lastIdent := self.getLastToken(0)
				self.pushToken("=", Equals)
				self.pushToken(lastIdent.Value, lastIdent.TokenType)
			}

			self.pushToken(currChar, reserved)
			continue
		}

		// --- Handle multi-character tokens ---

		if isDigit(self.at()) {
			self.tokenizeInt()
			continue
		}

		if isLetter(self.at()) {
			self.tokenizeIdentifier()
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

		return self.tokens, fmt.Errorf("%s:%d:%d: Unrecognized character '%s' found", config.Filename, self.currLn, self.currCol, self.at())
	}

	self.pushToken("EOF", EndOfFile)

	return self.tokens, nil
}

func (self *Lexer) tokenizeComment() {
	self.eat()
	comment := ""
	for self.isNotEof() && self.at() != "\n" {
		comment += self.eat()
	}
	self.pushToken(strings.TrimSpace(comment), Comment)
}

func (self *Lexer) tokenizeString() {
	self.eat()
	content := ""
	for self.isNotEof() && self.at() != "\"" {
		if self.at() == "\\" {
			content += self.eat()
		}
		content += self.eat()
	}
	self.eat()
	self.pushToken(content, Str)
}

func (self *Lexer) tokenizeInt() {
	num := ""
	for self.isNotEof() && isDigit(self.at()) {
		num += self.eat()
	}
	self.pushToken(num, Int)
}

func (self *Lexer) tokenizeIdentifier() {
	ident := ""
	for self.isNotEof() && (isLetter(self.at()) || isDigit(self.at())) {
		ident += self.eat()
	}

	// Check for reserved keywords
	if reserved, ok := keywords[ident]; ok {
		self.pushToken(ident, reserved)
	} else {
		self.pushToken(ident, Identifier)
	}
}

func (self *Lexer) isNotEof() bool {
	return len(self.sourceChars) > 0
}

func (self *Lexer) at() string {
	return self.sourceChars[0]
}

func (self *Lexer) next(offset int) string {
	if len(self.sourceChars) < offset+2 {
		return ""
	}
	return self.sourceChars[offset+1]
}

func (self *Lexer) eat() string {
	if len(self.sourceChars) == 0 {
		return ""
	}

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

func (self *Lexer) getLastToken(offset int) *Token {
	return self.tokens[len(self.tokens)-1-offset]
}

func isLetter(sourceChar string) bool {
	return unicode.IsLetter([]rune(sourceChar)[0])
}

func isDigit(sourceChar string) bool {
	return unicode.IsDigit([]rune(sourceChar)[0])
}

func isSkippable(sourceChar string) bool {
	return slices.Contains([]string{"\n", "\r", "\t", " "}, sourceChar)
}
